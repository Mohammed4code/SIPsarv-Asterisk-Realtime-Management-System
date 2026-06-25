package asterisk

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"sipsarv/internal/config"
)

func CheckConnection() bool {
	cfg := config.Load()
	host, port, _, _ := cfg.Asterisk.Get()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 5*time.Second)
	if err != nil {
		log.Printf("❌ فشل الاتصال بـ AMI: %v", err)
		return false
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	log.Printf("📡 AMI Banner: %s", strings.TrimSpace(line))
	return strings.Contains(line, "Asterisk")
}

func ReloadPJSIP() {
	conn, reader, err := AMISession()
	if err != nil {
		log.Printf("⚠️ تحذير: فشل إعادة تحميل PJSIP: %v", err)
		return
	}
	defer func() {
		conn.Write([]byte("Action: Logoff\r\n\r\n"))
		conn.Close()
	}()

	resp, err := sendAMICommand(conn, reader, "Action: Command\r\nCommand: pjsip reload\r\n\r\n")
	if err != nil {
		log.Printf("⚠️ تحذير: فشل إرسال أمر reload: %v", err)
		return
	}
	if strings.Contains(resp, "Success") {
		log.Println("🔄 PJSIP reload تم بنجاح")
	} else {
		log.Printf("⚠️ رد reload: %s", strings.TrimSpace(resp))
	}
}

func AMISession() (net.Conn, *bufio.Reader, error) {
	cfg := config.Load()
	host, port, user, password := cfg.Asterisk.Get()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 10*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("لا يمكن الاتصال بـ AMI: %v", err)
	}

	conn.SetDeadline(time.Now().Add(15 * time.Second))
	reader := bufio.NewReader(conn)

	// قراءة الترحيب
	if _, err := reader.ReadString('\n'); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("فشل قراءة ترحيب AMI: %v", err)
	}

	// تسجيل الدخول
	loginCmd := fmt.Sprintf("Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", user, password)
	if _, err := conn.Write([]byte(loginCmd)); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("فشل إرسال Login: %v", err)
	}

	resp, err := readAMIResponse(reader)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("فشل قراءة رد Login: %v", err)
	}
	if !strings.Contains(resp, "Success") {
		conn.Close()
		return nil, nil, fmt.Errorf("رُفض تسجيل الدخول: %s", resp)
	}

	log.Println("✅ تم تسجيل الدخول إلى AMI بنجاح")
	return conn, reader, nil
}

func sendAMICommand(conn net.Conn, reader *bufio.Reader, cmd string) (string, error) {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return "", err
	}
	return readAMIResponse(reader)
}

func readAMIResponse(reader *bufio.Reader) (string, error) {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		var block strings.Builder
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return block.String(), err
			}
			block.WriteString(line)
			if strings.TrimRight(line, "\r\n") == "" {
				break
			}
		}
		msg := block.String()
		if strings.HasPrefix(msg, "Event:") && !strings.Contains(msg, "Response:") {
			continue
		}
		return msg, nil
	}
	return "", fmt.Errorf("انتهت مهلة انتظار رد AMI")
}