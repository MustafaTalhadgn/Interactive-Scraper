# TorAnaliz - Cyber Threat Intelligence (CTI) Platform

Bu proje, Dark Web (Tor Ağı) üzerindeki siber tehdit kaynaklarını izleyen, toplanan verileri işleyip anlamlandıran ve görselleştiren uçtan uca bir Siber Tehdit İstihbaratı platformudur.


## 📖 Proje Özeti
TorAnaliz; Ransomware blogları, yasadışı forumlar ve sızıntı siteleri gibi .onion uzantılı kaynakları belirli aralıklarla tarar. Elde edilen ham verileri temizler, analiz eder ve bir analistin kolayca yorumlayabileceği şekilde Dashboard üzerinde görselleştirir. Sistem, verileri sadece depolamakla kalmaz, içerik analizi yaparak her veriye bir Kritiklik Skoru atar.

## 🛠 Kullanılan Teknolojiler
Proje, mikroservis mimarisine uygun olarak aşağıdaki teknolojilerle geliştirilmiştir:

* **Backend:** Go  - Gin Framework, GoQuery
* **Frontend:** React.js, TailwindCSS
* **Veritabanı:** PostgreSQL
* **Altyapı:** Docker & Docker Compose
* **Ağ:** Tor Proxy - peterdavehello/tor-socks-proxy

## 🚀 Kurulum ve Çalıştırma (Docker)
Proje tamamen Dockerize edilmiştir ve tek komutla ayağa kalkacak şekilde tasarlanmıştır.

1.  **Projeyi klonlayın:**
    ```
    git clone https://github.com/MustafaTalhadgn/Interactive-Scraper
    cd toranaliz
    ```

2.  **Sistemi başlatın:**
    Aşağıdaki komut, tüm servisleri (API, DB, Scraper, dashboard, Tor) derler ve başlatır:
    
    ```
    docker-compose up -d --build
    ```
    Dashboardı ayağa kaldırın 

    ```  
    npm run dev
    ```

3.  **Erişim:**

    Sistem ayağa kalktıktan sonra tarayıcınızdan erişebilirsiniz:

    * **Dashboard :** `http://localhost:5173/`
    * **Giriş:**      `Hesap oluşturarak sisteme girebilirsiniz`

> **Not:** İlk çalıştırmada Tor bağlantısının kurulması  ağ hızına bağlı olarak 30-60 saniye sürebilir.
> **Not:**  Veriler direk dashboarda düşmeyebilir scraperın taramasını bekleyin.
> **Not:** Eğer bilgisayarınızda 5173 portu doluysa 5174 portunu deneyebilirsiniz 

## 🧠 Başlık Üretim ve Analiz Mantığı
Sistem, Dark Web kaynaklarından veri çekerken aşağıdaki mantığı izler:

1.  **Veri Temizleme:** Kaynaktan çekilen HTML verisi, okunabilirliği artırmak için CSS ve Scriptlerden arındırılarak **Markdown** formatına çevrilir.
2.  **Başlık Oluşturma:**
    * Sistem öncelikle kaynaktaki `<h1>`, `<h2>` veya `article-title` etiketlerini arar.
    * Eğer net bir başlık bulunamazsa, içerik metninin ilk cümlesi veya özet kısmı "Otomatik Başlık" olarak atanır.
3.  **Kritiklik Puanlama (Scoring):** İçerik metni taranarak belirli anahtar kelimeler aranır (Örn: "database leak", "ransomware", "ssn", "password"). Bulunan her kritik kelime ve varlık (Bitcoin adresi vb.), içeriğin "Kritiklik Skoru"nu (0-100 arası) artırır.

## ⚠️ Yasal Uyarı
Bu proje sadece eğitim ve akademik araştırma amacıyla geliştirilmiştir. Proje kapsamında kaynaklardan hiçbir zararlı dosya, veritabanı dökümü veya yasadışı içerik **indirilmemekte ve saklanmamaktadır**. 