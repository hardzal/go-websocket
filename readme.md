# Go-websocket

- test with websocat

---

Level 1: Read data from websocket

Level 2: Intermediate

- Server websocket yang mana menyimpan banyak koneksi ke client
- Jika 1 client mengirim pesan, server akan broadcast ke semua client lainnya (Termasuk pengirim) (chat app?)
- tambahkan nama pengguna ke setiap client
- format broadcast: [username]: message
- gunakan goroutine, channel, sync.Mutex

Level 3: Room System

- Implementasikan system chat room
- Manajemen map[string][]\*Client
- JSON parsing

Level 4: Push Notification

- Ketika klient mengirim "start" -> server mengirim notifikasi tiap 3 detik
- jika klien mengirim "stop" -> server berhenti kirim notifikasi ke client itu

Level 5: WebSocket + REST hybrid

- buat 1 endpoint websocket /ws da 1 endpoint /notify
- Ketika ada POST /notify (dari curl /postman) semua client
