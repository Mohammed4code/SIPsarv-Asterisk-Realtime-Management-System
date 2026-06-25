
        // ---------- الإعدادات ----------
        const API_BASE = '/api';
        const ASTERISK_HOST = '192.168.8.24';

        // ---------- العناصر ----------
        const toast = document.getElementById('toast');
        const toastMessage = document.getElementById('toastMessage');
        const contactsContainer = document.getElementById('contactsContainer');

        // ---------- دوال عرض الرسائل ----------
        function showToast(message, type = 'success') {
            toast.className = 'toast show ' + type;
            toastMessage.textContent = message;
            setTimeout(function() { toast.className = 'toast'; }, 4000);
        }

        // ---------- جلب البيانات من النموذج ----------
        function getFormData() {
            return {
                full_name: document.getElementById('fullName').value.trim(),
                passport: document.getElementById('passport').value.trim(),
                address: document.getElementById('address').value.trim(),
                phone_number: document.getElementById('phoneNumber').value.trim(),
                pin: document.getElementById('pinCode').value.trim(),
                sip_username: document.getElementById('sipUsername').value.trim(),
                sip_secret: document.getElementById('sipSecret').value.trim(),
                email: document.getElementById('email').value.trim(),
                connection_type: document.getElementById('connectionType').value,
                reg_timeout: parseInt(document.getElementById('regTimeout').value) || 120,
                asterisk_notes: document.getElementById('asteriskNotes').value.trim()
            };
        }

        // ---------- التحقق من البيانات ----------
        function validateData(data) {
            if (!data.full_name) { showToast('⚠️ الاسم الكامل مطلوب', 'error'); return false; }
            if (!data.phone_number) { showToast('⚠️ رقم الهاتف مطلوب', 'error'); return false; }
            if (!data.pin) { showToast('⚠️ PIN مطلوب', 'error'); return false; }
            if (!data.sip_username) { showToast('⚠️ اسم مستخدم SIP مطلوب', 'error'); return false; }
            return true;
        }

        // ---------- إرسال طلب إلى Go ----------
        async function sendToGo(endpoint, data) {
            try {
                showToast('⏳ جاري الإرسال...', 'loading');

                const response = await fetch(API_BASE + endpoint, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json'
                    },
                    body: JSON.stringify(data)
                });

                const result = await response.json();

                if (!response.ok) {
                    throw new Error(result.message || 'HTTP ' + response.status);
                }

                showToast(result.message || '✅ تمت العملية بنجاح', 'success');
                loadContacts();
                return result;
            } catch (error) {
                console.error('❌ Error:', error);
                showToast('❌ فشل: ' + error.message, 'error');
                return null;
            }
        }

        // ---------- تحميل قائمة جهات الاتصال ----------
        async function loadContacts() {
            try {
                const response = await fetch(API_BASE + '/contacts/list');
                const result = await response.json();

                if (result.success && result.data && result.data.contacts) {
                    const contacts = result.data.contacts;
                    if (contacts.length === 0) {
                        contactsContainer.innerHTML = '<div class="no-contacts">📭 لا توجد جهات اتصال محفوظة</div>';
                        return;
                    }

                    let html = `
                        <table class="contacts-table">
                            <thead>
                                <tr>
                                    <th>#</th>
                                    <th>الاسم</th>
                                    <th>رقم الهاتف</th>
                                    <th>SIP Username</th>
                                    <th>الإجراء</th>
                                </tr>
                            </thead>
                            <tbody>
                    `;

                    contacts.forEach(function(contact, index) {
                        html += `
                            <tr>
                                <td>${index + 1}</td>
                                <td>${contact.full_name || '-'}</td>
                                <td>${contact.phone_number || '-'}</td>
                                <td>${contact.sip_username || '-'}</td>
                                <td>
                                    <button class="delete-btn" onclick="deleteContact('${contact.phone_number}')" title="حذف">
                                        <i class="fas fa-trash-alt"></i>
                                    </button>
                                </td>
                            </tr>
                        `;
                    });

                    html += '</tbody></table>';
                    contactsContainer.innerHTML = html;
                }
            } catch (error) {
                console.error('❌ فشل تحميل جهات الاتصال:', error);
            }
        }

        // ---------- حذف جهة اتصال ----------
        async function deleteContact(phoneNumber) {
            if (!confirm('هل أنت متأكد من حذف جهة الاتصال رقم ' + phoneNumber + '؟')) return;

            try {
                const response = await fetch(API_BASE + '/contacts/delete', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ phone_number: phoneNumber })
                });

                const result = await response.json();
                if (result.success) {
                    showToast(result.message, 'success');
                    loadContacts();
                } else {
                    showToast(result.message, 'error');
                }
            } catch (error) {
                showToast('❌ فشل الحذف: ' + error.message, 'error');
            }
        }

        // ---------- تحميل إعدادات Asterisk ----------
        async function loadAsteriskConfig() {
            try {
                const response = await fetch(API_BASE + '/asterisk/config');
                const result = await response.json();
                
                if (result.success && result.data) {
                    document.getElementById('asteriskHost').value = result.data.host || '';
                    document.getElementById('asteriskPort').value = result.data.port || '';
                    document.getElementById('amiUser').value = result.data.user || '';
                }
            } catch (error) {
                console.error('❌ فشل تحميل إعدادات Asterisk:', error);
            }
        }

        // ---------- تحديث إعدادات Asterisk ----------
        async function updateAsteriskConfig() {
            const host = document.getElementById('asteriskHost').value.trim();
            const port = document.getElementById('asteriskPort').value.trim();
            const user = document.getElementById('amiUser').value.trim();
            const password = document.getElementById('amiPassword').value.trim();
            
            if (!host) {
                showToast('⚠️ عنوان الخادم مطلوب', 'error');
                return;
            }
            
            const data = {
                host: host,
                port: port || '5038',
                user: user || 'admin'
            };
            
            if (password) {
                data.password = password;
            }
            
            try {
                showToast('⏳ جاري تحديث الإعدادات...', 'loading');
                
                const response = await fetch(API_BASE + '/asterisk/config', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json'
                    },
                    body: JSON.stringify(data)
                });
                
                const result = await response.json();
                
                if (result.success) {
                    showToast(result.message, 'success');
                    setTimeout(checkServerStatus, 500);
                } else {
                    showToast(result.message || '❌ فشل التحديث', 'error');
                }
            } catch (error) {
                console.error('❌ Error updating config:', error);
                showToast('❌ فشل التحديث: ' + error.message, 'error');
            }
        }

        // ---------- فحص حالة الخادم ----------
        async function checkServerStatus() {
            const dot = document.getElementById('serverDot');
            const statusText = document.getElementById('serverStatusText');

            dot.className = 'dot checking';
            statusText.textContent = '⏳ جاري الفحص...';

            try {
                // تحميل الإعدادات أولاً للحصول على العنوان الحالي
                const configResponse = await fetch(API_BASE + '/asterisk/config');
                const configResult = await configResponse.json();
                
                if (configResult.success && configResult.data) {
                    const host = configResult.data.host || 'غير معروف';
                    document.getElementById('asteriskHost').value = host;
                    document.getElementById('asteriskPort').value = configResult.data.port || '5038';
                    document.getElementById('amiUser').value = configResult.data.user || 'admin';
                }

                // فحص Go Server
                const healthResponse = await fetch(API_BASE + '/health');
                if (healthResponse.ok) {
                    dot.className = 'dot';
                    dot.style.background = '#3dbb7a';
                    dot.style.boxShadow = '0 0 8px #3dbb7a';
                    statusText.textContent = '✅ متصل بـ Go Server';
                } else {
                    throw new Error('Go server not healthy');
                }

                // فحص Asterisk
                const asteriskResponse = await fetch(API_BASE + '/asterisk/status');
                const asteriskResult = await asteriskResponse.json();
                const badge = document.getElementById('asteriskStatusBadge');

                if (asteriskResult.success && asteriskResult.data && asteriskResult.data.connected) {
                    const host = asteriskResult.data.host || 'غير معروف';
                    badge.className = 'badge badge-status';
                    badge.innerHTML = '<i class="fas fa-circle"></i> ✅ متصل بـ ' + host;
                } else {
                    const host = asteriskResult.data && asteriskResult.data.host ? asteriskResult.data.host : 'غير معروف';
                    badge.className = 'badge badge-status offline';
                    badge.innerHTML = '<i class="fas fa-circle"></i> ❌ غير متصل بـ ' + host;
                }

            } catch (error) {
                console.error('❌ Health check error:', error);
                dot.className = 'dot offline';
                dot.style.background = '#b33c3c';
                dot.style.boxShadow = '0 0 8px #b33c3c';
                statusText.textContent = '❌ غير متصل بـ Go Server';

                const badge = document.getElementById('asteriskStatusBadge');
                badge.className = 'badge badge-status offline';
                badge.innerHTML = '<i class="fas fa-circle"></i> ❌ غير متصل';
            }
        }

        // ---------- ربط الأزرار ----------
        document.getElementById('addToAsteriskBtn').addEventListener('click', async function() {
            const data = getFormData();
            if (!validateData(data)) return;
            await sendToGo('/extensions/add', data);
        });

        document.getElementById('saveContactBtn').addEventListener('click', async function() {
            const data = getFormData();
            if (!validateData(data)) return;
            await sendToGo('/contacts/save', data);
        });

        document.getElementById('refreshStatus').addEventListener('click', checkServerStatus);
        document.getElementById('updateConfigBtn').addEventListener('click', updateAsteriskConfig);

        // إعادة تعيين
        document.querySelector('button[type="reset"]').addEventListener('click', function() {
            toast.className = 'toast';
        });

        // منع إرسال النموذج
        document.getElementById('contactForm').addEventListener('submit', function(e) {
            e.preventDefault();
        });

        // ---------- التحميل الأولي ----------
        setTimeout(function() {
            loadAsteriskConfig();
            checkServerStatus();
            loadContacts();
        }, 500);

        // ---------- نافذة الامتدادات ----------
        const extensionsOverlay = document.getElementById('extensionsOverlay');
        const extensionsContainer = document.getElementById('extensionsContainer');
        const extensionsCount = document.getElementById('extensionsCount');

        async function loadExtensions() {
            extensionsContainer.innerHTML = '<div class="no-contacts">⏳ جاري التحميل...</div>';
            try {
                const response = await fetch(API_BASE + '/extensions/list');
                const result = await response.json();

                if (!result.success) {
                    extensionsContainer.innerHTML = '<div class="no-contacts">❌ ' + (result.message || 'فشل التحميل') + '</div>';
                    extensionsCount.textContent = '';
                    return;
                }

                const extensions = (result.data && result.data.extensions) || [];
                extensionsCount.textContent = extensions.length + ' امتداد';

                if (extensions.length === 0) {
                    extensionsContainer.innerHTML = '<div class="no-contacts">📭 لا توجد امتدادات مسجلة حالياً</div>';
                    return;
                }

                let html = `
                    <table class="contacts-table">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>الامتداد</th>
                                <th>الاسم</th>
                                <th>رقم الهاتف</th>
                                <th>Transport</th>
                                <th>الإجراء</th>
                            </tr>
                        </thead>
                        <tbody>
                `;

                extensions.forEach(function(ext, index) {
                    html += `
                        <tr>
                            <td>${index + 1}</td>
                            <td><strong>${ext.id}</strong></td>
                            <td>${ext.full_name || '-'}</td>
                            <td>${ext.phone_number || '-'}</td>
                            <td>${ext.transport || '-'}</td>
                            <td>
                                <button class="delete-btn" onclick="deleteExtension('${ext.id}')" title="حذف">
                                    <i class="fas fa-trash-alt"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });

                html += '</tbody></table>';
                extensionsContainer.innerHTML = html;
            } catch (error) {
                console.error('❌ فشل تحميل الامتدادات:', error);
                extensionsContainer.innerHTML = '<div class="no-contacts">❌ فشل الاتصال بالخادم</div>';
            }
        }

        async function deleteExtension(id) {
            if (!confirm('هل أنت متأكد من حذف الامتداد ' + id + '؟')) return;
            try {
                const response = await fetch(API_BASE + '/extensions/delete', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ id: id })
                });
                const result = await response.json();
                showToast(result.message, result.success ? 'success' : 'error');
                if (result.success) loadExtensions();
            } catch (error) {
                showToast('❌ فشل الحذف: ' + error.message, 'error');
            }
        }

        document.getElementById('showExtensionsBtn').addEventListener('click', function() {
            extensionsOverlay.classList.add('show');
            loadExtensions();
        });

        document.getElementById('closeExtensionsBtn').addEventListener('click', function() {
            extensionsOverlay.classList.remove('show');
        });

        document.getElementById('refreshExtensionsBtn').addEventListener('click', loadExtensions);

        extensionsOverlay.addEventListener('click', function(e) {
            if (e.target === extensionsOverlay) {
                extensionsOverlay.classList.remove('show');
            }
        });

        // ---------- تحديث الحالة كل 30 ثانية ----------
        setInterval(checkServerStatus, 30000);