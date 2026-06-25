package asterisk

import (
	
	"fmt"
	"strings"

	"sipsarv/internal/config"
	"sipsarv/internal/database"
)

func InsertExtensionRealtime(c database.Contact) error {
	authID := "auth" + c.SipUsername

	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("فشل بدء transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// التحقق من وجود الامتداد
	var exists int
	tx.QueryRow("SELECT COUNT(*) FROM ps_endpoints WHERE id = ?", c.SipUsername).Scan(&exists)
	if exists > 0 {
		return fmt.Errorf("الامتداد %s موجود مسبقاً في قاعدة البيانات", c.SipUsername)
	}

	// إضافة Auth
	_, err = tx.Exec(`
		INSERT INTO ps_auths (id, auth_type, username, password)
		VALUES (?, 'userpass', ?, ?)
	`, authID, c.SipUsername, c.SipSecret)
	if err != nil {
		return fmt.Errorf("فشل إضافة Auth: %v", err)
	}

	// إضافة AOR
	_, err = tx.Exec(`
		INSERT INTO ps_aors (id, max_contacts)
		VALUES (?, 3)
	`, c.SipUsername)
	if err != nil {
		return fmt.Errorf("فشل إضافة AOR: %v", err)
	}

	// تحديد النقل
	transport := "transport-ws"
	switch strings.ToLower(c.ConnectionType) {
	case "tcp":
		transport = "transport-tcp"
	case "udp":
		transport = "transport-udp"
	}

	cfg := config.Load()
	_, _, _, asteriskHost := cfg.Asterisk.Get()

	// إضافة Endpoint
	_, err = tx.Exec(`
		INSERT INTO ps_endpoints
			(id, transport, aors, auth, context, disallow, allow, callerid, 
			 webrtc, dtls_auto_generate_cert, media_address)
		VALUES
			(?, ?, ?, ?, 'default', 'all', 'ulaw,alaw', ?, 'yes', 'yes', ?)
	`, c.SipUsername, transport, c.SipUsername, authID,
		fmt.Sprintf("%s <%s>", c.FullName, c.PhoneNumber), asteriskHost)
	if err != nil {
		return fmt.Errorf("فشل إضافة Endpoint: %v", err)
	}

	// إضافة بيانات الموظف
	regTimeout := c.RegTimeout
	if regTimeout <= 0 {
		regTimeout = 3600
	}
	_, err = tx.Exec(`
		INSERT INTO contacts
			(full_name, passport, address, phone_number, pin, sip_username, sip_secret,
			 email, connection_type, reg_timeout, asterisk_notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, c.FullName, c.Passport, c.Address, c.PhoneNumber, c.Pin, c.SipUsername, c.SipSecret,
		c.Email, c.ConnectionType, regTimeout, c.AsteriskNotes)
	if err != nil {
		return fmt.Errorf("فشل حفظ بيانات الموظف: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("فشل تأكيد العملية: %v", err)
	}

	return nil
}

func DeleteExtensionRealtime(id string) error {
	authID := "auth" + id

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	// ترتيب الحذف مهم بسبب القيد FK
	tx.Exec("DELETE FROM contacts WHERE sip_username = ?", id)
	tx.Exec("DELETE FROM ps_endpoints WHERE id = ?", id)
	tx.Exec("DELETE FROM ps_aors WHERE id = ?", id)
	tx.Exec("DELETE FROM ps_auths WHERE id = ?", authID)

	return tx.Commit()
}

func ListExtensions() ([]database.Extension, error) {
	rows, err := database.DB.Query(`
		SELECT e.id, e.transport, e.aors, e.auth, e.context, e.disallow, e.allow, 
		       COALESCE(a.max_contacts, 0),
		       COALESCE(c.full_name, ''),
		       COALESCE(c.phone_number, '')
		FROM ps_endpoints e
		LEFT JOIN ps_aors a ON a.id = e.aors
		LEFT JOIN contacts c ON c.sip_username = e.id
		ORDER BY e.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	extensions := []database.Extension{}
	for rows.Next() {
		var ext database.Extension
		err := rows.Scan(&ext.ID, &ext.Transport, &ext.Aors, &ext.Auth, &ext.Context,
			&ext.Disallow, &ext.Allow, &ext.MaxContact, &ext.FullName, &ext.Phone)
		if err != nil {
			continue
		}
		extensions = append(extensions, ext)
	}
	return extensions, nil
}