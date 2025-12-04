# Spesifikasi Teknis Lengkap: Sistem Electronic Medical Record (EMR) Rumah Sakit

**Versi:** 1.0
**Tanggal:** 2025-11-06
**Penulis:** MiniMax Agent

---

## 1. Ringkasan Eksekutif

Dokumen ini menyajikan spesifikasi teknis lengkap untuk pengembangan dan implementasi sistem Electronic Medical Record (EMR) generasi baru yang dirancang untuk meningkatkan efisiensi operasional, perawatan pasien, dan kepatuhan regulasi di lingkungan rumah sakit modern. Sistem ini dibangun di atas tumpukan teknologi modern yang mengutamakan skalabilitas, keamanan, dan interoperabilitas, terdiri dari **backend microservices dengan Golang**, **frontend web responsif dengan React.js & Tailwind CSS**, dan **database cloud-native Neon PostgreSQL**.

**Tujuan Utama Proyek:**

1.  **Digitalisasi Komprehensif:** Menggantikan proses manual dan sistem warisan dengan platform digital terpusat untuk semua data klinis dan administratif pasien.
2.  **Peningkatan Kualitas Perawatan:** Memberikan akses real-time kepada tenaga medis terhadap riwayat pasien yang lengkap dan akurat, mendukung pengambilan keputusan klinis yang lebih baik.
3.  **Efisiensi Operasional:** Mengotomatiskan alur kerja klinis dan administratif, mulai dari pendaftaran pasien, penjadwalan, rekam medis, hingga penagihan dan pelaporan.
4.  **Interoperabilitas Data:** Memfasilitasi pertukaran data yang mulus dengan sistem internal lainnya (ERP, Laboratorium, Radiologi) dan eksternal (HIE, BPJS) melalui standar industri seperti HL7 FHIR dan DICOM.
5.  **Kepatuhan Regulasi:** Memastikan kepatuhan penuh terhadap standar keamanan dan privasi data kesehatan nasional (UU PDP, PMK 24/2022) dan internasional (HIPAA).

**Value Proposition:**

Sistem EMR ini bukan sekadar pengganti kertas, melainkan sebuah platform strategis yang memungkinkan rumah sakit untuk bertransformasi secara digital. Dengan arsitektur microservices yang fleksibel, sistem ini dapat beradaptasi dengan kebutuhan masa depan, mengintegrasikan teknologi baru seperti AI dan analitik data, serta memberikan fondasi yang kokoh untuk model perawatan berbasis nilai (value-based care). Integrasi yang erat dengan sistem Enterprise Resource Planning (ERP) akan menyatukan data klinis dan finansial, memberikan pandangan 360 derajat terhadap operasional rumah sakit dan mendorong efisiensi biaya.

---

## 2. Kebutuhan dan Spesifikasi Sistem

Bagian ini merinci kebutuhan fungsional dan non-fungsional yang harus dipenuhi oleh sistem EMR.

### 2.1. Kebutuhan Fungsional (Functional Requirements)

Kebutuhan fungsional dikelompokkan berdasarkan modul inti sistem.

| Modul | Fitur | Deskripsi | Prioritas |
| :--- | :--- | :--- | :--- |
| **Manajemen Pasien** | Pendaftaran Pasien (Admission) | Pendaftaran pasien baru dan lama, pengumpulan data demografis, asuransi, dan kontak darurat. Sistem harus mendukung pencarian pasien duplikat. | Wajib |
| | Manajemen Janji Temu | Penjadwalan, pembatalan, dan penjadwalan ulang janji temu dengan dokter. Integrasi dengan kalender dokter dan pengiriman notifikasi/pengingat otomatis. | Wajib |
| | Alur Pasien (Patient Flow) | Pelacakan status dan lokasi pasien secara real-time di dalam rumah sakit (mis., Pendaftaran → Poliklinik → Laboratorium → Farmasi → Kasir). | Wajib |
| **Rekam Medis Elektronik** | Catatan Klinis (Clinical Notes) | Pembuatan dan pengelolaan catatan klinis terstruktur (SOAP) dan tidak terstruktur. Dukungan template untuk berbagai spesialisasi. | Wajib |
| | Daftar Masalah (Problem List) | Pencatatan diagnosis pasien yang aktif dan lampau menggunakan standar kode internasional (ICD-10). | Wajib |
| | Riwayat Pengobatan (Medication History) | Pencatatan riwayat alergi dan semua obat yang pernah dan sedang dikonsumsi pasien. | Wajib |
| | E-Prescription (Resep Elektronik) | Pembuatan resep digital yang terintegrasi dengan modul farmasi dan apotek. Mendukung standar EPCS (Electronic Prescription for Controlled Substances). | Wajib |
| | Manajemen Dokumen | Unggah dan kelola dokumen medis pasien (hasil lab, scan, foto, dll.) dalam berbagai format. | Wajib |
| **Integrasi Klinis** | Integrasi Laboratorium (LIS) | Permintaan tes lab (Order Entry) dan penerimaan hasil secara elektronik langsung ke dalam EMR pasien. | Wajib |
| | Integrasi Radiologi (RIS/PACS) | Permintaan pemeriksaan radiologi dan akses ke gambar medis (mis., X-ray, CT Scan) melalui viewer DICOM terintegrasi (DICOMweb). | Wajib |
| **Penagihan & Pelaporan**| Encounter & Billing | Pencatatan semua layanan yang diberikan selama encounter pasien untuk proses penagihan. Integrasi dengan sistem ERP untuk sinkronisasi data tagihan. | Wajib |
| | Pelaporan Kepatuhan | Generasi laporan untuk standar kepatuhan seperti UDS (Uniform Data System) 2024, eCQM (electronic Clinical Quality Measures), dan Promoting Interoperability (PI). | Wajib |
| | Dasbor Analitik | Dasbor visual untuk manajemen dan tenaga medis yang menampilkan Key Performance Indicators (KPIs) rumah sakit (mis., waktu tunggu, BOR). | Tinggi |
| **Akses Pengguna** | Dasbor Berbasis Peran (Role-Based) | Tampilan antarmuka yang disesuaikan untuk setiap peran pengguna (Dokter, Perawat, Resepsionis, Admin, Pasien) yang hanya menampilkan informasi dan fungsi yang relevan. | Wajib |
| | Akses Mobile | Aplikasi web yang sepenuhnya responsif dan dapat diakses dari perangkat mobile (tablet, smartphone) dengan mematuhi standar aksesibilitas WCAG 2.1 AA. | Wajib |
| | Portal Pasien | Portal aman bagi pasien untuk melihat ringkasan rekam medis mereka, jadwal janji temu, hasil lab, dan berkomunikasi dengan tenaga medis. | Tinggi |

### 2.2. Kebutuhan Non-Fungsional (Non-Functional Requirements)

| Kategori | Kebutuhan | Spesifikasi |
| :--- | :--- | :--- |
| **Performa** | Waktu Respons API | Latensi P95 untuk semua endpoint API inti harus < 200 ms. |
| | Waktu Muat Halaman | Waktu muat halaman (Load Time) untuk dasbor utama tidak boleh melebihi 3 detik. |
| | Konkurensi | Sistem harus mampu menangani 1.000 pengguna konkuren dengan performa yang terjaga. |
| **Keamanan** | Enkripsi Data | Semua data (at-rest dan in-transit) harus dienkripsi menggunakan standar AES-256 dan TLS 1.3. |
| | Kontrol Akses | Implementasi Role-Based Access Control (RBAC) yang ketat dan Multi-Factor Authentication (MFA) berbasis FIDO2/WebAuthn. |
| | Audit Trail | Semua akses dan perubahan data pada data medis sensitif (PHI) harus dicatat dalam log audit yang tidak dapat diubah (immutable) dan disimpan selama minimal 25 tahun. |
| **Kepatuhan (Compliance)**| Kepatuhan Regulasi | Sistem harus mematuhi **HIPAA Security Rule** dan **UU PDP Indonesia (UU No. 27 Tahun 2022)**, termasuk PMK 24/2022 tentang retensi data EMR. |
| | Interoperabilitas | Kepatuhan penuh pada standar **HL7 FHIR R4** untuk pertukaran data kesehatan dan **DICOMweb** untuk pencitraan medis. |
| **Skalabilitas**| Arsitektur | Arsitektur microservices harus memungkinkan penskalaan independen untuk setiap layanan (mis., layanan otentikasi, layanan pasien) berdasarkan beban kerja. |
| | Database | Database harus mendukung penskalaan komputasi dan penyimpanan secara dinamis tanpa downtime. |
| **Ketersediaan (Availability)**| Uptime Sistem | Sistem harus memiliki ketersediaan minimal 99.9% (uptime). |
| | Disaster Recovery | Harus ada rencana pemulihan bencana (Disaster Recovery Plan) dengan target RPO (Recovery Point Objective) < 1 jam dan RTO (Recovery Time Objective) < 4 jam. |
| **Pemeliharaan** | Deployment | Proses deployment dan rollback harus otomatis (CI/CD) dan dapat dilakukan tanpa downtime (zero-downtime deployment). |
| | Monitoring | Sistem harus memiliki observability penuh dengan structured logging, distributed tracing, dan monitoring metrik real-time. |

---

## 3. Gambaran Umum Arsitektur

Arsitektur sistem EMR ini dirancang dengan pendekatan modern, cloud-native, dan berorientasi layanan (service-oriented) untuk memenuhi kebutuhan skalabilitas, keamanan, dan fleksibilitas.

### 3.1. Arsitektur Tingkat Tinggi (High-Level Architecture)

Arsitektur sistem mengadopsi model **Microservices** yang didekomposisi berdasarkan Domain-Driven Design (DDD). Setiap layanan (service) memiliki tanggung jawab yang jelas dan berkomunikasi satu sama lain melalui API yang terdefinisi dengan baik.

![Figure 1: High-Level Microservices Architecture](https://i.imgur.com/example.png)  *(Catatan: Diagram ini adalah representasi konseptual. Diagram teknis yang lebih detail akan disediakan di bagian arsitektur backend.)*

Komponen utama arsitektur ini adalah:

1.  **Clients (Frontend):** Aplikasi web berbasis React.js yang diakses oleh pengguna melalui browser. Didesain dengan pendekatan **Backend-for-Frontend (BFF)**, di mana ada lapisan API yang khusus melayani kebutuhan antarmuka pengguna, mengoptimalkan data yang dikirim ke client.
2.  **API Gateway:** Pintu masuk tunggal untuk semua permintaan dari client. Bertanggung jawab untuk routing, otentikasi (verifikasi JWT), rate limiting, dan logging. Komponen ini menyederhanakan interaksi client dan melindungi layanan backend.
3.  **Backend Services (Golang):** Kumpulan microservices independen yang dibangun dengan Golang. Setiap layanan mengelola domain bisnis tertentu (mis., `Patient Service`, `Encounter Service`, `Auth Service`). Mereka berkomunikasi secara sinkron (melalui REST/gRPC) untuk permintaan langsung dan asinkron (melalui Message Broker) untuk proses latar belakang.
4.  **FHIR-Native Data Plane:** Sebuah lapisan data terpusat yang berfungsi sebagai "single source of truth" untuk data klinis, diimplementasikan menggunakan **FHIR Server**. Ini memastikan konsistensi data dan kepatuhan standar interoperabilitas. Microservices berinteraksi dengan lapisan ini untuk mengelola data rekam medis.
5.  **Database (Neon PostgreSQL):** Database relasional serverless yang menyimpan data aplikasi, termasuk data non-klinis (user, role, audit log) dan data yang mendukung operasional microservices.
6.  **Message Broker (NATS/Kafka):** Memfasilitasi komunikasi asinkron dan event-driven antar microservices. Digunakan untuk proses seperti notifikasi, sinkronisasi data ke ERP, dan alur kerja yang berjalan lama (Saga pattern).
7.  **Sistem Eksternal:** Termasuk sistem ERP rumah sakit, sistem laboratorium (LIS), sistem radiologi (RIS/PACS), dan platform Health Information Exchange (HIE) yang terintegrasi melalui `Integration Service`.

### 3.2. Arsitektur Database

Database utama menggunakan **Neon PostgreSQL**, platform database serverless yang memisahkan komputasi dan penyimpanan, memungkinkan skalabilitas elastis dan efisiensi biaya. Desain skema database didasarkan pada Entity-Relationship Diagram (ERD) berikut, yang mencakup domain inti EMR.

![Figure 2: Entity-Relationship Diagram (ERD) untuk Database EMR.](/workspace/docs/design/emr_database_erd.png)

**Prinsip Desain Database:**

*   **Normalisasi:** Skema dinormalisasi untuk mengurangi redundansi data dan memastikan integritas.
*   **Partisi Data:** Tabel besar seperti `audit_logs` dan `encounters` akan dipartisi berdasarkan rentang waktu (mis., bulanan) untuk meningkatkan performa query dan mempermudah manajemen data.
*   **Keamanan:** Implementasi **Row-Level Security (RLS)** untuk memastikan tenaga medis hanya dapat mengakses data pasien yang menjadi tanggung jawab mereka. Enkripsi at-rest dan in-transit diaktifkan secara default.
*   **Sinkronisasi ERP:** Menggunakan **Change Data Capture (CDC)** untuk menangkap perubahan data (mis., data tagihan baru) dan mempublikasikannya ke message broker untuk sinkronisasi near real-time dengan sistem ERP.

### 3.3. Arsitektur Keamanan (Zero Trust)

Sistem ini mengadopsi model keamanan **Zero Trust**, di mana tidak ada kepercayaan implisit yang diberikan kepada entitas mana pun di dalam atau di luar jaringan.

*   **Autentikasi Kuat:** Semua akses, baik oleh pengguna maupun layanan, harus diautentikasi secara ketat menggunakan MFA dan token (JWT).
*   **Enkripsi End-to-End:** Komunikasi antar layanan dienkripsi menggunakan **mutual TLS (mTLS)** yang di-enforce oleh service mesh (seperti Istio atau Linkerd).
*   **Least Privilege Access:** Setiap pengguna dan layanan hanya diberikan izin minimum yang diperlukan untuk menjalankan fungsinya.
*   **Micro-segmentation:** Jaringan disegmentasi untuk mengisolasi layanan satu sama lain, membatasi dampak jika terjadi kompromi pada salah satu layanan.

Pendekatan arsitektur ini memastikan bahwa sistem tidak hanya fungsional tetapi juga aman, tangguh, dan siap untuk masa depan.

---

## 4. Spesifikasi Backend (Golang)

Backend sistem EMR dibangun menggunakan bahasa pemrograman **Golang** dengan arsitektur **microservices**. Pilihan ini didasarkan pada performa tinggi Golang, konkurensi bawaan (goroutines), keamanan (strong typing, memory safety), dan ekosistem yang matang untuk aplikasi cloud-native.

### 4.1. Dekomposisi Layanan (Service Decomposition)

Layanan dipecah berdasarkan kapabilitas bisnis (Bounded Contexts dari DDD) untuk memastikan otonomi dan kohesi yang tinggi:

| Service Name | Bounded Context | Responsibilities |
| :--- | :--- | :--- |
| **Auth Service** | Identity & Access | Manajemen pengguna, autentikasi (login, JWT generation), otorisasi (RBAC), manajemen sesi, dan SSO. |
| **Patient Service** | Patient Administration | Mengelola data demografis pasien, pendaftaran, pencarian, dan manajemen consent. |
| **Encounter Service** | Clinical Encounter | Mengelola semua interaksi klinis, dari admisi hingga discharge, termasuk diagnosis (ICD-10) dan catatan klinis. |
| **Scheduling Service** | Appointment Management | Mengelola jadwal janji temu dokter, ketersediaan, pemesanan, dan notifikasi pengingat. |
| **Terminology Service** | Clinical Terminology | Menyediakan akses terpusat ke standar terminologi medis seperti ICD-10, SNOMED CT, dan LOINC. |
| **Orders Service** | Clinical Orders | Mengelola permintaan untuk tes laboratorium, pemeriksaan radiologi, dan resep obat. |
| **Results Service** | Clinical Results | Menerima dan menyimpan hasil tes laboratorium dan laporan radiologi dari sistem eksternal (LIS/RIS). |
| **Integration Service** | External Systems | Bertindak sebagai fasilitator untuk integrasi dengan sistem eksternal seperti ERP, LIS, RIS, dan HIE. Mengelola transformasi data (mis., HL7v2 ke FHIR). |
| **Notification Service** | User Communication | Mengelola pengiriman notifikasi real-time (WebSocket), email, dan SMS kepada pengguna (pasien dan staf). |

### 4.2. Desain dan Spesifikasi API

*   **Prinsip Desain:** Semua API dirancang dengan pendekatan **RESTful** dan **FHIR-native**.
*   **Versioning:** Versi API akan dicantumkan di URL (mis., `/api/v1/patients`).
*   **Format Data:** JSON akan menjadi format standar untuk semua request dan response.
*   **Dokumentasi:** Spesifikasi API akan didokumentasikan menggunakan **OpenAPI 3.1.0** dan disajikan melalui Swagger UI untuk kemudahan pengujian dan integrasi.
*   **Keamanan:** Semua endpoint akan dilindungi dan memerlukan token JWT (Bearer Token) yang valid, yang akan diverifikasi di API Gateway.
*   **Contoh Endpoint (Patient Service):**
    *   `GET /api/v1/patients`: Mendapatkan daftar pasien (dengan paginasi).
    *   `POST /api/v1/patients`: Mendaftarkan pasien baru.
    *   `GET /api/v1/patients/{id}`: Mendapatkan detail pasien berdasarkan ID.
    *   `PUT /api/v1/patients/{id}`: Memperbarui data demografis pasien.
    *   `GET /api/v1/patients/{id}/encounters`: Mendapatkan riwayat semua pertemuan klinis untuk pasien tertentu.

### 4.3. Komunikasi Antar Layanan

*   **Komunikasi Sinkron:** Untuk permintaan yang memerlukan respons langsung (mis., frontend mengambil data pasien), layanan akan berkomunikasi melalui **REST API** atau **gRPC** untuk efisiensi yang lebih tinggi.
*   **Komunikasi Asinkron:** Untuk proses yang tidak memerlukan respons segera, tahan lama, atau perlu di-decouple (mis., mengirim notifikasi, sinkronisasi ke ERP, saga-pattern untuk pemesanan lab), layanan akan menggunakan **Message Broker** (NATS atau Kafka). Ini meningkatkan resiliensi sistem; jika layanan penerima sedang down, pesan akan tetap ada di antrian dan diproses nanti.

### 4.4. Framework dan Library

*   **Web Framework:** Dipertimbangkan menggunakan framework populer seperti **Gin**, **Echo**, atau **Fiber** karena performa tinggi dan middleware yang kaya.
*   **ORM/Database Library:** `database/sql` standar dengan driver `pgx` untuk PostgreSQL. GORM atau SQLC dapat digunakan untuk mempercepat pengembangan.
*   **Keamanan:** Library standar `crypto` untuk enkripsi, `golang-jwt/jwt` untuk manajemen JWT, dan `go-webauthn` untuk implementasi FIDO2.

---

## 5. Desain Database (Neon PostgreSQL)

Database merupakan fondasi sistem, dan **Neon PostgreSQL** dipilih karena arsitektur serverless-nya yang menawarkan skalabilitas on-demand, fitur branching untuk development, dan kepatuhan HIPAA.

### 5.1. Desain Skema dan Tabel

Desain skema (seperti pada ERD di Bagian 3.2) mencakup tabel-tabel utama berikut:

*   `users`, `roles`, `permissions`: Untuk sistem RBAC.
*   `patients`, `patient_demographics`: Menyimpan data inti pasien.
*   `encounters`, `clinical_notes`, `diagnoses`: Untuk mencatat setiap interaksi klinis.
*   `appointments`: Untuk modul penjadwalan.
*   `lab_orders`, `lab_results`: Untuk integrasi LIS.
*   `audit_trails`: Tabel krusial untuk mencatat semua akses dan modifikasi PHI.

### 5.2. Keamanan dan Kepatuhan Database

*   **Enkripsi:** Neon menyediakan enkripsi **AES-256 at-rest** secara default. Koneksi ke database akan di-enforce menggunakan **TLS 1.3 (verify-full)**.
*   **Row-Level Security (RLS):** Kebijakan RLS akan diimplementasikan secara ekstensif. Contoh:
    *   Seorang dokter hanya bisa melihat data pasien yang memiliki janji temu dengannya atau berada di bawah perawatannya.
    *   Pasien (melalui portal) hanya dapat melihat data medis miliknya sendiri.
    *   *Script SQL untuk kebijakan ini terdapat di `docs/design/rls_security_policies.sql`.*
*   **Audit Trail:** Menggunakan ekstensi **pgAudit**, semua query (SELECT, INSERT, UPDATE, DELETE) pada tabel yang berisi PHI akan dicatat secara otomatis ke tabel `audit_trails`.
*   **Manajemen Akses:** Akses ke database dari aplikasi akan menggunakan kredensial yang dirotasi secara berkala dan disimpan dengan aman di secret manager (mis., HashiCorp Vault atau AWS Secrets Manager).

### 5.3. Kinerja dan Skalabilitas

*   **Penskalaan Otomatis:** Kemampuan serverless Neon memungkinkan compute resources untuk mati saat tidak digunakan (scale to zero) dan menyala kembali dalam hitungan detik, sangat efisien untuk lingkungan development dan staging.
*   **Read Replicas:** Untuk beban kerja read-heavy seperti analitik dan pelaporan, **Read Replicas** akan digunakan untuk memisahkan beban kerja dari database utama (primary).
*   **Indexing Strategy:** Indeks akan dibuat secara strategis pada kolom yang sering digunakan dalam query `WHERE`, `JOIN`, dan `ORDER BY`, terutama pada foreign keys dan kolom pencarian seperti `patient_mrn` (Medical Record Number).

### 5.4. Backup dan Disaster Recovery

*   **Point-in-Time Recovery (PITR):** Neon menyediakan fitur PITR yang memungkinkan restore database ke kondisi kapan pun dalam periode retensi yang ditentukan (mis., 30 hari untuk plan Scale). Ini krusial untuk pemulihan dari kesalahan operasional.
*   **Disaster Recovery (DR):** Meskipun Neon mengelola replikasi di dalam satu region, strategi DR antar-region akan dirancang dengan melakukan backup logis (menggunakan `pg_dump`) secara berkala ke cloud storage di region yang berbeda.

---

## 6. Spesifikasi Frontend (React + Tailwind CSS)

Antarmuka pengguna (UI) adalah wajah sistem. Dibangun dengan **React.js** untuk arsitektur berbasis komponen yang dinamis dan **Tailwind CSS** untuk utility-first styling yang cepat dan konsisten.

### 6.1. Arsitektur Komponen

*   **Atomic Design:** Mengadopsi metodologi Atomic Design untuk menstrukturkan komponen:
    *   **Atoms:** Komponen paling dasar (Button, Input, Label).
    *   **Molecules:** Kombinasi dari atoms (Search bar = Input + Button).
    *   **Organisms:** Bagian UI yang lebih kompleks (Form pendaftaran, Header navigasi).
    *   **Templates:** Layout halaman tanpa data.
    *   **Pages:** Template yang diisi dengan data nyata.
*   **Component Library:** Sebuah library komponen yang dapat digunakan kembali akan dikembangkan (mis., menggunakan Storybook) untuk memastikan konsistensi visual dan mempercepat pengembangan.

### 6.2. Manajemen State (State Management)

*   **Local State:** `useState` dan `useReducer` dari React untuk mengelola state internal komponen.
*   **Global State:** Untuk state yang perlu dibagikan ke seluruh aplikasi (mis., informasi pengguna yang login, role), akan digunakan **Zustand** atau **Redux Toolkit**.
*   **Server Cache State:** Untuk mengelola data dari API (fetching, caching, updating), akan digunakan **React Query** atau **SWR**. Ini akan menangani caching data, re-fetching di background, dan optimisme UI-updates secara otomatis.

### 6.3. UI/UX dan Desain Responsif

*   **Desain Responsif:** Dengan Tailwind CSS, semua komponen dan layout akan dirancang dengan pendekatan **mobile-first**, memastikan pengalaman pengguna yang optimal di berbagai ukuran layar, dari smartphone hingga desktop.
*   **Aksesibilitas:** Kepatuhan pada **Web Content Accessibility Guidelines (WCAG) 2.1 Level AA** adalah wajib. Ini termasuk penggunaan atribut ARIA yang benar, kontras warna yang memadai, navigasi keyboard, dan label untuk semua elemen interaktif.
*   **Alur Pengguna (User Flow):** Desain akan berfokus pada penyederhanaan alur kerja klinis yang kompleks. Misalnya, proses input data SOAP akan dibuat seintuitif mungkin dengan template dan auto-suggestion.

### 6.4. Interaksi dengan Backend

*   **BFF (Backend-for-Frontend):** Frontend tidak akan langsung memanggil setiap microservice. Sebaliknya, ia akan berkomunikasi dengan **API Gateway** atau lapisan BFF yang akan mengagregasi data dari berbagai layanan. Ini menyederhanakan logika di sisi frontend dan mengurangi jumlah network call.
*   **Real-time Updates:** Untuk fitur real-time (mis., notifikasi hasil lab baru, update status antrian), frontend akan membuka koneksi **WebSocket** ke **Notification Service**.
---

## 7. Spesifikasi Integrasi ERP

Integrasi antara sistem EMR dan Enterprise Resource Planning (ERP) rumah sakit adalah kunci untuk menyatukan data klinis dan operasional, menciptakan efisiensi, dan mendukung pengambilan keputusan berbasis data di seluruh organisasi.

### 7.1. Arsitektur Integrasi

Arsitektur integrasi akan bersifat **API-first** dan **event-driven**, menggunakan **Integration Service** sebagai middleware atau hub terpusat. Pendekatan ini menghindari integrasi point-to-point yang rapuh.

1.  **Event-Driven Architecture (EDA):** Sistem EMR akan mempublikasikan event bisnis penting (mis., `patient_discharged`, `lab_order_created`, `invoice_generated`) ke **Message Broker**. `Integration Service` akan men-subscribe event-event ini.
2.  **Integration Service:** Layanan ini bertanggung jawab untuk:
    *   **Transformasi Data:** Mengubah format data dari EMR (mis., FHIR) ke format yang dibutuhkan oleh ERP (mis., HL7v2, CSV, atau API proprietary ERP).
    *   **Orkestrasi Proses:** Mengelola alur kerja yang melibatkan kedua sistem.
    *   **Error Handling & Retry Logic:** Menangani kegagalan koneksi atau penolakan data dari ERP dengan mekanisme antrian dan coba lagi (retry).
3.  **API ERP:** `Integration Service` akan berinteraksi dengan API yang diekspos oleh sistem ERP untuk mengirim dan mengambil data.

### 7.2. Alur Data dan Titik Integrasi

Integrasi akan difokuskan pada tiga domain utama:

| Domain | Alur Kerja | Data yang Disinkronkan | Arah | Pola Integrasi |
| :--- | :--- | :--- | :--- | :--- |
| **Keuangan (Financial)** | **Encounter-to-Cash** | Data demografi pasien, diagnosis (ICD-10), layanan/tindakan medis yang diberikan, informasi asuransi. | EMR → ERP | Event-Driven (Real-time) |
| | **Rekonsiliasi Tagihan** | Status pembayaran, detail klaim asuransi, penyesuaian tagihan. | ERP → EMR | API Call (Batch/Scheduled) |
| **Inventaris (Inventory)** | **Manajemen Stok Farmasi** | Permintaan obat (resep) dari EMR akan mengurangi stok di ERP. Notifikasi stok rendah dari ERP ke EMR. | EMR ↔ ERP | Event-Driven & API Call |
| | **Manajemen Aset Medis** | Penggunaan alat medis habis pakai selama prosedur (dicatat di EMR) akan mengurangi inventaris di ERP. | EMR → ERP | Event-Driven |
| **Sumber Daya Manusia (HR)** | **Manajemen Staf Klinis** | Sinkronisasi data master penyedia layanan kesehatan (dokter, perawat) dari ERP ke EMR, termasuk jadwal, kredensial, dan status aktif. | ERP → EMR | API Call (Batch/Scheduled) |

### 7.3. Standar dan Teknologi

*   **HL7 FHIR:** Akan menjadi standar utama untuk representasi data klinis di sisi EMR.
*   **HL7 v2:** Untuk kompatibilitas dengan sistem ERP warisan, `Integration Service` akan mampu memproses dan menghasilkan pesan HL7 v2 (mis., ADT untuk pendaftaran, ORM untuk order, ORU untuk hasil).
*   **REST/SOAP API:** Interaksi dengan API modern atau legacy dari sistem ERP.

---

## 8. Keamanan & Kepatuhan

Keamanan dan kepatuhan adalah aspek non-negotiable dalam sistem EMR. Sistem ini dirancang dari awal dengan prinsip **Security by Design** dan untuk memenuhi regulasi yang ketat.

### 8.1. Kepatuhan Regulasi

*   **HIPAA (Health Insurance Portability and Accountability Act):**
    *   **Security Rule:** Implementasi penuh terhadap *Administrative*, *Physical*, dan *Technical Safeguards*. Ini termasuk enkripsi (8.2), kontrol akses (8.3), audit (8.4), dan Rencana Kontinjensi (DRP).
    *   **Privacy Rule:** Fitur sistem akan memastikan hak pasien terpenuhi, seperti hak untuk mengakses dan meminta amendemen pada PHI mereka.
    *   **Business Associate Agreement (BAA):** BAA akan ditandatangani dengan semua vendor pihak ketiga yang menangani PHI, termasuk penyedia cloud (AWS/GCP) dan Neon.
*   **UU PDP & Regulasi Indonesia:**
    *   **UU No. 27 Tahun 2022 (PDP):** Sistem akan mengelola persetujuan (consent) pasien secara eksplisit untuk pemrosesan data, serta menyediakan mekanisme untuk permintaan penghapusan data sesuai hak subjek data.
    *   **PMK No. 24 Tahun 2022:** Log audit dan data rekam medis akan disimpan selama minimal **25 tahun**, sesuai dengan persyaratan retensi untuk EMR.

### 8.2. Enkripsi dan Perlindungan Data

*   **Encryption in Transit:** Semua komunikasi (client-server, server-server) akan diwajibkan menggunakan **TLS 1.3**.
*   **Encryption at Rest:** Semua data PHI di database dan penyimpanan file akan dienkripsi menggunakan **AES-256**.
*   **Manajemen Kunci:** Kunci enkripsi akan dikelola menggunakan layanan Key Management Service (KMS) dari penyedia cloud untuk memisahkan kunci dari data.

### 8.3. Manajemen Identitas dan Akses (IAM)

*   **Autentikasi:** **Multi-Factor Authentication (MFA)** akan diwajibkan untuk semua pengguna, terutama yang memiliki akses ke PHI. Metode **FIDO2/WebAuthn** akan diutamakan karena ketahanannya terhadap phishing.
*   **Otorisasi:** **Role-Based Access Control (RBAC)** akan diimplementasikan secara granular. Izin tidak hanya akan terikat pada peran (mis., 'Dokter') tetapi juga pada konteks (mis., 'Dokter yang merawat pasien X').
*   **Prinsip Least Privilege:** Pengguna hanya akan memiliki akses ke data dan fungsi yang mutlak diperlukan untuk pekerjaan mereka.

### 8.4. Audit dan Monitoring

*   **Comprehensive Audit Trails:** Setiap tindakan yang terkait dengan PHI (Create, Read, Update, Delete, Print, Export) akan dicatat. Log akan berisi **SIAPA** (user ID), **APA** (tindakan), **KAPAN** (timestamp), dan **DARI MANA** (IP address/workstation).
*   **Immutability:** Log audit akan disimpan di media WORM (Write-Once, Read-Many) atau menggunakan teknologi blockchain/ledger untuk memastikan tidak dapat diubah.
*   **Real-time Monitoring:** Sistem akan dipantau secara aktif untuk aktivitas mencurigakan (mis., akses data dalam jumlah besar di luar jam kerja) menggunakan alat SIEM (Security Information and Event Management).

---

## 9. Fitur Real-time

Kemampuan untuk menyajikan informasi secara real-time sangat penting dalam lingkungan klinis yang dinamis. Sistem ini akan mengimplementasikan beberapa fitur real-time menggunakan teknologi **WebSocket** dan **Event Streaming**.

### 9.1. Arsitektur Notifikasi Real-time

1.  **Notification Service (Golang):** Sebuah microservice khusus yang mengelola koneksi WebSocket dari semua klien yang aktif.
2.  **Message Broker:** Ketika sebuah event terjadi di sistem (mis., `Results Service` menerima hasil lab baru), event tersebut dipublikasikan ke message broker.
3.  **WebSocket Push:** `Notification Service` men-subscribe event-event relevan dari broker. Ketika sebuah event diterima, layanan ini akan mengidentifikasi pengguna mana yang perlu diberi tahu dan mengirimkan pesan ke klien yang sesuai melalui koneksi WebSocket yang ada.

### 9.2. Kasus Penggunaan

*   **Notifikasi Hasil Kritis:** Dokter menerima notifikasi instan di layar mereka ketika hasil laboratorium kritis untuk pasien mereka tersedia.
*   **Update Antrian Pasien:** Layar di ruang tunggu dan dasbor perawat diperbarui secara real-time ketika status pasien berubah (mis., 'dipanggil', 'sedang di dalam ruangan', 'selesai').
*   **Chat Kolaborasi:** Memungkinkan staf medis untuk berkomunikasi secara aman dan real-time mengenai perawatan pasien, menggantikan aplikasi perpesanan konsumen yang tidak aman.
*   **Live Dashboard:** KPI di dasbor manajemen (mis., jumlah pasien saat ini, waktu tunggu rata-rata) diperbarui secara live.

### 9.3. Implementasi Teknis

*   **Backend:** Menggunakan library WebSocket untuk Golang (mis., `gorilla/websocket`). `Notification Service` akan menjaga state koneksi aktif dan pemetaan user-to-connection.
*   **Frontend:** Klien React akan menggunakan API WebSocket bawaan browser atau library seperti `Socket.IO-client` untuk membuat dan mendengarkan koneksi.
*   **Skalabilitas:** `Notification Service` akan dirancang untuk menjadi stateless sebisa mungkin, memungkinkan beberapa instance berjalan di belakang load balancer untuk menangani puluhan ribu koneksi konkuren.
---

## 10. Deployment & Infrastruktur

Infrastruktur sistem EMR akan dirancang untuk menjadi **cloud-native**, memanfaatkan platform cloud modern (seperti AWS, GCP, atau Azure) untuk mencapai skalabilitas, keandalan, dan keamanan yang tinggi.

### 10.1. Infrastruktur Cloud

*   **Containerization:** Semua microservices backend akan di-package sebagai **Docker containers**. Ini memastikan konsistensi lingkungan dari development hingga produksi.
*   **Orchestration:** **Kubernetes (K8s)** akan digunakan sebagai platform orkestrasi kontainer. K8s akan mengelola deployment, penskalaan (scaling), dan self-healing dari layanan secara otomatis.
*   **Infrastructure as Code (IaC):** Seluruh infrastruktur (VPC, cluster K8s, database, load balancer) akan didefinisikan sebagai kode menggunakan **Terraform**. Ini memungkinkan provisi dan modifikasi infrastruktur yang dapat direproduksi dan diaudit.
*   **Service Mesh:** **Istio** atau **Linkerd** akan diimplementasikan di dalam cluster K8s untuk menyediakan kapabilitas krusial seperti mTLS (enkripsi antar layanan), load balancing cerdas, circuit breaking, dan observability lalu lintas layanan.

### 10.2. Pipeline CI/CD

Pipeline Continuous Integration/Continuous Deployment (CI/CD) yang sepenuhnya otomatis adalah kunci untuk rilis yang cepat dan andal. Alat seperti **GitLab CI**, **GitHub Actions**, atau **Jenkins** akan digunakan.

**Alur Pipeline:**

1.  **Commit:** Developer melakukan `git push` ke repositori.
2.  **Build & Test (CI):**
    *   Pipeline CI terpicu secara otomatis.
    *   Menjalankan `go build` untuk kompilasi.
    *   Menjalankan unit tests dan static code analysis.
    *   Membangun Docker image dan mendorongnya ke container registry (mis., ECR, GCR).
3.  **Deploy to Staging (CD):**
    *   Setelah CI berhasil, image baru secara otomatis di-deploy ke lingkungan **Staging**.
    *   Menjalankan integration tests dan end-to-end tests terhadap lingkungan Staging.
4.  **Deploy to Production:**
    *   Deployment ke **Produksi** memerlukan persetujuan manual (manual approval gate).
    *   Menggunakan strategi deployment **Blue-Green** atau **Canary Release** untuk meminimalkan risiko.
        *   **Blue-Green:** Deploy versi baru ke lingkungan siaga (Green). Setelah verifikasi, lalu lintas dialihkan dari versi lama (Blue) ke Green. Memungkinkan rollback instan dengan mengalihkan lalu lintas kembali ke Blue.
        *   **Canary:** Rilis versi baru ke sebagian kecil pengguna (mis., 5%). Jika tidak ada masalah, lalu lintas secara bertahap ditingkatkan hingga 100%.

### 10.3. Monitoring dan Observability

Observability penuh adalah wajib untuk mendeteksi dan mendiagnosis masalah secara proaktif. Tiga pilar observability akan diimplementasikan:

*   **Logging:** Semua layanan akan menghasilkan **structured logs** (dalam format JSON). Log ini akan dikumpulkan, di-parse, dan diindeks menggunakan tumpukan **Fluentd/Filebeat + Elasticsearch + Kibana (EFK/ELK Stack)**.
*   **Metrics:** Metrik performa (latency, throughput, error rate) dan metrik sumber daya (CPU, memori) akan diekspos oleh setiap layanan menggunakan format **Prometheus**. Prometheus akan mengumpulkan metrik ini, dan **Grafana** akan digunakan untuk membuat dasbor visualisasi dan alert.
*   **Tracing:** **Distributed Tracing** akan diimplementasikan menggunakan **OpenTelemetry**. Ini memungkinkan pelacakan sebuah permintaan saat melewati berbagai microservices, sangat berguna untuk mengidentifikasi bottleneck performa dalam sistem yang terdistribusi.

---

## 11. Strategi Pengujian (Testing Strategy)

Pengujian adalah proses berlapis yang dirancang untuk memastikan kualitas, keamanan, dan kepatuhan sistem dari level terendah hingga sistem secara keseluruhan.

### 11.1. Piramida Pengujian

Strategi pengujian akan mengikuti model piramida pengujian:

1.  **Unit Tests (Dasar Piramida):**
    *   **Tujuan:** Memverifikasi fungsionalitas dari unit kode terkecil (fungsi/metode) secara terisolasi.
    *   **Framework:** Menggunakan library testing bawaan Golang (`testing`) dan library assertion seperti `testify/assert`.
    *   **Cakupan:** Setiap layanan harus memiliki cakupan unit test > 80%.

2.  **Integration Tests (Tengah Piramida):**
    *   **Tujuan:** Memverifikasi interaksi antar komponen, terutama antara layanan dan databasenya, atau antar dua layanan.
    *   **Contoh:** Menguji apakah `Patient Service` dapat membuat entri di database PostgreSQL dengan benar, atau apakah `Encounter Service` dapat berhasil memanggil `Patient Service`.
    *   **Tools:** Menggunakan Docker untuk menjalankan instance database atau layanan dependen dalam lingkungan pengujian.

3.  **End-to-End (E2E) Tests / UI Tests (Puncak Piramida):**
    *   **Tujuan:** Mensimulasikan alur kerja pengguna nyata dari awal hingga akhir, dari antarmuka pengguna hingga ke database dan kembali lagi.
    *   **Contoh:** Mensimulasikan alur "seorang dokter login, mencari pasien, menambahkan diagnosis, dan meresepkan obat".
    *   **Framework:** **Cypress** atau **Playwright** untuk mengotomatiskan interaksi browser.

### 11.2. Pengujian Non-Fungsional

*   **Performance Testing:**
    *   **Load Testing:** Menggunakan alat seperti **k6** atau **JMeter** untuk mensimulasikan ribuan pengguna dan memastikan sistem memenuhi target waktu respons dan throughput.
    *   **Stress Testing:** Mendorong sistem hingga ke batasnya untuk mengidentifikasi titik kegagalan dan memastikan sistem dapat pulih dengan baik (graceful degradation).
*   **Security Testing:**
    *   **Static Application Security Testing (SAST):** Analisis kode sumber untuk menemukan kerentanan umum (diintegrasikan dalam pipeline CI).
    *   **Dynamic Application Security Testing (DAST):** Memindai aplikasi yang sedang berjalan untuk menemukan kerentanan seperti SQL Injection, XSS, dll.
    *   **Penetration Testing:** Melibatkan pihak ketiga (ethical hacker) untuk mencoba menembus sistem sebelum dan sesudah go-live.

### 11.3. Pengujian Kepatuhan (Compliance Testing)

*   **Tujuan:** Memastikan semua kontrol teknis untuk HIPAA dan regulasi lainnya diimplementasikan dan berfungsi dengan benar.
*   **Contoh Skenario:**
    *   Memverifikasi bahwa log audit dibuat untuk setiap akses PHI.
    *   Menguji kebijakan RBAC/RLS untuk memastikan seorang dokter tidak dapat melihat data pasien di luar wewenangnya.
    *   Memastikan data yang diekspor dienkripsi.
---

## 12. Peta Jalan Implementasi (Implementation Roadmap)

Peta jalan ini menguraikan implementasi sistem EMR dalam fase-fase yang dapat dikelola, memungkinkan pengiriman nilai secara bertahap dan iteratif.

### Fase 1: Fondasi dan Core Services (Bulan 1-3)

Fokus pada pembangunan fondasi arsitektur dan fungsionalitas inti EMR.

*   **Milestones:**
    *   **Bulan 1: Fondasi Arsitektur & Keamanan.**
        *   Setup infrastruktur cloud (Kubernetes, VPC, Neon DB) menggunakan Terraform.
        *   Implementasi pipeline CI/CD dasar.
        *   Inisialisasi `Auth Service` dengan RBAC dan `User Service`.
        *   Desain final skema database dan ERD.
    *   **Bulan 2: Manajemen Pasien & Penjadwalan.**
        *   Pengembangan `Patient Service` (pendaftaran, pencarian).
        *   Pengembangan `Scheduling Service` (pembuatan janji temu).
        *   Frontend: Halaman login, dasbor, modul pendaftaran pasien, dan kalender janji temu.
    *   **Bulan 3: Rekam Medis Dasar.**
        *   Pengembangan `Encounter Service` (pencatatan encounter, diagnosis).
        *   Pengembangan `Terminology Service`.
        *   Frontend: Halaman rekam medis pasien (menampilkan demografi & riwayat encounter).
*   **Deliverables:** MVP (Minimum Viable Product) internal dengan fungsionalitas pendaftaran, penjadwalan, dan pencatatan diagnosis dasar. Siap untuk UAT (User Acceptance Testing) internal oleh tim super-user.

### Fase 2: Integrasi Klinis dan Fitur Lanjutan (Bulan 4-6)

Memperluas kapabilitas klinis sistem dan mengintegrasikannya dengan sistem lain.

*   **Milestones:**
    *   **Bulan 4: E-Prescription & Orders.**
        *   Pengembangan `Orders Service` untuk resep elektronik.
        *   Integrasi awal dengan sistem farmasi (jika API tersedia).
        *   Frontend: Modul E-Prescription.
    *   **Bulan 5: Integrasi Laboratorium & Radiologi.**
        *   Pengembangan `Integration Service` untuk memproses pesan HL7 dari LIS/RIS.
        *   Pengembangan `Results Service` untuk menyimpan hasil lab/radiologi.
        *   Frontend: Tab hasil lab dan radiologi di rekam medis pasien, termasuk viewer DICOM dasar.
    *   **Bulan 6: Fitur Real-time & Pelaporan Awal.**
        *   Pengembangan `Notification Service` (WebSocket).
        *   Implementasi notifikasi real-time untuk hasil lab kritis.
        *   Pengembangan modul pelaporan dasar (mis., jumlah pasien harian).
*   **Deliverables:** Sistem EMR fungsional yang terintegrasi dengan LIS/RIS. Siap untuk pilot terbatas di satu atau dua departemen.

### Fase 3: Integrasi ERP dan Go-Live (Bulan 7-9)

Fokus pada sinkronisasi dengan ERP dan persiapan untuk peluncuran skala penuh.

*   **Milestones:**
    *   **Bulan 7: Integrasi Keuangan (ERP).**
        *   Pengembangan alur `Encounter-to-Cash` di `Integration Service` untuk mengirim data tagihan ke ERP.
        *   Frontend: Finalisasi modul penagihan.
    *   **Bulan 8: Migrasi Data & Pelatihan.**
        *   Pengembangan skrip untuk migrasi data demografi pasien dari sistem lama.
        *   Pelaksanaan sesi pelatihan komprehensif untuk semua staf (dokter, perawat, admin).
    *   **Bulan 9: Go-Live & Hypercare.**
        *   Peluncuran sistem EMR di seluruh rumah sakit.
        *   Periode "Hypercare": Dukungan intensif di tempat selama 2-4 minggu pertama setelah go-live untuk mengatasi masalah dengan cepat.
*   **Deliverables:** Sistem EMR yang beroperasi penuh di seluruh rumah sakit, dengan dukungan pasca-implementasi.

### Fase 4: Optimisasi dan Pengembangan Lanjutan (Bulan 10+)

*   **Milestones:**
    *   Pengembangan Portal Pasien.
    *   Pengembangan dasbor analitik lanjutan untuk manajemen.
    *   Optimisasi performa berdasarkan data penggunaan nyata.
    *   Implementasi fitur-fitur baru berdasarkan umpan balik pengguna.

---

## 13. Penilaian dan Mitigasi Risiko

Identifikasi proaktif terhadap potensi risiko sangat penting untuk keberhasilan proyek.

| Kategori Risiko | Deskripsi Risiko | Dampak | Probabilitas | Strategi Mitigasi |
| :--- | :--- | :--- | :--- | :--- |
| **Teknis** | **Keterlambatan Pengembangan** | Proyek melebihi jadwal dan anggaran. | Sedang | Mengadopsi metodologi Agile (Scrum) dengan sprint 2-mingguan untuk memantau kemajuan secara ketat. Mengidentifikasi dan mengatasi penghambat (blockers) setiap hari. |
| | **Masalah Skalabilitas** | Performa sistem menurun drastis saat beban pengguna tinggi setelah go-live. | Rendah | Melakukan load testing yang ketat sebelum go-live. Menggunakan arsitektur microservices dan database serverless yang memungkinkan penskalaan elastis. |
| | **Kesulitan Integrasi** | API dari sistem warisan (ERP/LIS) tidak terdokumentasi dengan baik atau tidak stabil. | Tinggi | Melakukan sesi penemuan (discovery) teknis yang mendalam dengan vendor sistem pihak ketiga di awal proyek. Mengembangkan adapter/konektor yang kuat dengan logging dan mekanisme retry yang andal. |
| **Kepatuhan** | **Pelanggaran Data (Data Breach)** | Kerusakan reputasi, denda regulasi (HIPAA/PDP), dan hilangnya kepercayaan pasien. | Rendah | Menerapkan arsitektur keamanan Zero Trust, enkripsi end-to-end, MFA, dan RBAC yang ketat. Melakukan penetration testing secara berkala. |
| | **Kegagalan Audit Kepatuhan** | Gagal menunjukkan kepatuhan pada regulator. | Rendah | Mengimplementasikan logging audit yang komprehensif dan tidak dapat diubah (immutable) sejak hari pertama. Mengotomatiskan pengujian kepatuhan dalam pipeline CI/CD. |
| **Operasional** | **Penolakan Pengguna (User Adoption)** | Staf medis menolak menggunakan sistem baru dan kembali ke cara lama, membuat investasi sia-sia. | Tinggi | Melibatkan perwakilan pengguna akhir (dokter, perawat) dalam proses desain (User-Centered Design). Menyediakan pelatihan yang ekstensif dan dukungan pasca-go-live (Hypercare). Menunjuk "champion users" di setiap departemen. |
| | **Migrasi Data Gagal** | Data pasien dari sistem lama tidak akurat atau tidak lengkap setelah migrasi. | Sedang | Melakukan beberapa kali uji coba migrasi (dry-run) ke lingkungan staging. Melakukan validasi dan pembersihan data sebelum migrasi final. |

---

## 14. Apendiks

Dokumen ini didukung oleh serangkaian artefak desain dan penelitian yang lebih rinci, yang dapat ditemukan di direktori `docs/` proyek. Ini termasuk:

*   **A. Dokumentasi API:** Spesifikasi OpenAPI 3.1.0 lengkap untuk semua microservices (tersedia di `docs/design/api_endpoints.md`).
*   **B. Skema Database:** Skrip DDL (Data Definition Language) SQL untuk membuat semua tabel, indeks, dan relasi (tersedia di `docs/design/sql_ddl_scripts.sql`).
*   **C. Kebijakan Keamanan Database:** Skrip SQL yang mendefinisikan kebijakan Row-Level Security (RLS) untuk berbagai peran pengguna (tersedia di `docs/design/rls_security_policies.sql`).
*   **D. Diagram Arsitektur:** Diagram ERD dan diagram arsitektur tingkat tinggi (tersedia di `docs/design/emr_database_erd.png` dan `docs/design/microservices_architecture.md`).
*   **E. Analisis Kepatuhan:** Analisis mendalam tentang pemetaan fitur sistem terhadap persyaratan HIPAA dan UU PDP (tersedia di `docs/compliance/medical_compliance_analysis.md`).