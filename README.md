# Sistem EMR Rumah Sakit - Backend

Sistem Rekam Medis Elektronik (EMR) komprehensif berbasis cloud-native yang dibangun dengan **arsitektur microservices Golang**. Dirancang untuk lingkungan rumah sakit modern, sistem ini memprioritaskan **keamanan Zero Trust**, **kepatuhan regulasi (HIPAA, UU PDP)**, dan **interoperabilitas tinggi**.

---

## 🎯 Ringkasan Eksekutif & Value Proposition

Proyek ini merupakan platform transformasi digital strategis bagi penyedia layanan kesehatan. Berbeda dengan EMR monolitik lama, sistem ini menawarkan:

*   **Skalabilitas**: Microservices yang terdekomposisi memungkinkan penskalaan independen untuk modul dengan beban tinggi (misalnya, Layanan Pasien) tanpa memengaruhi modul lain.
*   **Ketahanan (Resilience)**: Arsitektur event-driven memastikan bahwa kegagalan pada satu layanan (misalnya, Penagihan) tidak menghentikan operasi klinis kritis.
*   **Keamanan Utama**: Dibangun dari awal dengan pola pikir **Zero Trust**, menerapkan enkripsi saat istirahat/transit, RBAC yang ketat, dan jejak audit yang komprehensif.
*   **Interoperabilitas**: Dukungan asli untuk **HL7 FHIR R4** dan **DICOM** memastikan integrasi yang mulus dengan ekosistem kesehatan global.

---

## 🏗️ Arsitektur Sistem

Sistem ini mengikuti pendekatan **Domain-Driven Design (DDD)**, yang terstruktur menjadi microservices otonom.

### Arsitektur Tingkat Tinggi
*   **Backend-for-Frontend (BFF)**: API Gateway bertindak sebagai BFF, mengagregasi data untuk memberikan respons yang dioptimalkan bagi klien React, mengurangi lalu lintas jaringan.
*   **Event-Driven Backbone**: Menggunakan **NATS** untuk komunikasi asinkron. Alur kritis seperti "Pasien Pulang" memicu event yang dikonsumsi secara independen oleh layanan hilir (Penagihan, Notifikasi).
*   **Data Plane**: Lapisan data **FHIR-native** memastikan semua data klinis tersimpan dalam format standar, menjamin masa depan aset data rumah sakit.

### Microservices Inti
| Layanan | Tanggung Jawab | Teknologi Kunci |
| :--- | :--- | :--- |
| **Auth Service** | Identitas, penerbitan JWT, RBAC, MFA | `golang-jwt`, `webauthn` |
| **Patient Service** | Demografi, Pendaftaran, Indeks Pasien Utama | PostgreSQL, GORM |
| **Encounter Service** | Catatan klinis (SOAP), Diagnosis (ICD-10) | FHIR Stores |
| **Integration Service** | Transformasi HL7v2/FHIR untuk sistem ERP/LIS lama | Parser HL7 kustom |
| **Notification Service** | Peringatan real-time (Hasil lab, Update antrian) | WebSockets |

---

## 💡 Keputusan Teknis Kunci (Poin Bicara Interview)

### Mengapa Golang?
*   **Konkurensi**: Goroutine ringan Go sangat cocok untuk menangani ribuan koneksi WebSocket bersamaan untuk dasbor rumah sakit real-time.
*   **Performa**: Pengetikan statis dan kompilasi menawarkan performa mendekati C dengan keamanan memori, sangat penting untuk API kesehatan ber-throughput tinggi.

### Mengapa Neon PostgreSQL (Serverless)?
*   **Efisiensi Biaya**: Komputasi penskalaan otomatis berarti kita hanya membayar untuk penggunaan aktif, ideal untuk beban rumah sakit yang fluktuatif (hari sibuk vs malam sepi).
*   **Branching**: Fitur branching Neon memungkinkan pengembang untuk secara instan membuat klon database "copy-on-write" untuk menguji fitur tanpa memengaruhi data produksi.

### Mengapa NATS?
*   **Kesederhanaan & Kecepatan**: NATS menawarkan latensi dan kompleksitas operasional yang lebih rendah dibandingkan Kafka, sesuai dengan kebutuhan kami akan pesan internal berkecepatan tinggi.
*   **Decoupling**: Memungkinkan pola "Fire and Forget" untuk tugas non-kritis (misalnya, mengirim email konfirmasi), meningkatkan latensi yang dirasakan pengguna.

---

## 🚀 Fitur Utama

### Klinis & Operasional
- **Smart Encounters**: Catatan SOAP terstruktur dengan penyimpanan otomatis dan pencarian kode ICD-10.
- **E-Prescription**: Resep digital dengan pemeriksaan interaksi otomatis.
- **Lab & Radiologi**: Pemesanan terintegrasi dan peninjauan hasil (dukungan viewer DICOM).
- **Penjadwalan Sumber Daya**: Pemesanan bebas konflik untuk dokter dan ruangan.

### Keamanan & Kepatuhan (Krusial)
- **Log Audit Tidak Dapat Diubah**: Setiap tindakan baca/tulis pada Informasi Kesehatan Pasien (PHI) dicatat secara kriptografis selama 25 tahun (Kepatuhan: PMK 24/2022).
- **Row-Level Security (RLS)**: Kebijakan database memastikan dokter *hanya* dapat mengakses pasien yang mereka rawat.
- **Enkripsi**: AES-256 untuk penyimpanan database dan TLS 1.3 untuk semua lalu lintas API.

---

## 🛠️ Tech Stack

- **Bahasa**: Golang 1.21+
- **Framework**: Gin (HTTP), GORM (ORM)
- **Database**: Neon PostgreSQL
- **Messaging**: NATS
- **Caching**: Redis
- **Infrastruktur**: Docker, Kubernetes (Manifest disertakan)
- **Dokumentasi**: OpenAPI 3.1.0 (Swagger)

---

## 🏃‍♂️ Mulai Cepat (Dev Lokal)

### Prasyarat
- Go 1.21+
- Docker & Docker Compose

### Langkah-langkah
1.  **Clone & Setup**
    ```bash
    git clone https://github.com/hospital-emr/backend.git
    cd backend
    cp .env.example .env
    ```

2.  **Jalankan Infrastruktur (DB, Redis, NATS)**
    ```bash
    docker-compose up -d
    ```

3.  **Jalankan Migrasi & Seed Data**
    ```bash
    go run cmd/migrate/main.go up
    go run cmd/seed/main.go
    ```

4.  **Jalankan API**
    ```bash
    go run cmd/api/main.go
    # API tersedia di http://localhost:8080
    # Dokumen Swagger di http://localhost:8080/api/docs
    ```

---

## 🧪 Strategi Pengujian

*   **Unit Tests**: Validasi logika bisnis (`make test-unit`).
*   **Integration Tests**: Alur API-ke-Database (`make test-integration`).
*   **Security Scan**: Analisis statis untuk kerentanan (`make security-scan`).

---

## 👥 Tim & Peran

*   **Backend Architect**: Desain sistem, Implementasi keamanan.
*   **DevOps Engineer**: Pipeline CI/CD, Orkestrasi Kubernetes.
*   **Frontend Developer**: UI/UX React.js.
*   **QA Engineer**: Kepatuhan dan pengujian E2E.

---

> *Proyek ini adalah demonstrasi arsitektur sistem misi-kritis dengan kepatuhan tinggi yang cocok untuk lingkungan kesehatan perusahaan.*
 