# Foreign7 - Dil Öğrenme Uygulaması Görev Listesi

## ✅ Tamamlanan Özellikler
- [x] **Proje Kurulumu:** Go, Ollama SDK ve Bubble Tea entegrasyonu.
- [x] **Mimari (Refactoring):** `tearouter` ile stack-based navigasyon ve `internal/model` altında modüler yapı.
- [x] **Merkezi Veri Yönetimi:** `internal/data` ile `.txt` tabanlı çoklu liste ve aktif liste yönetimi.
- [x] **Mod 1 (Anlam Tahmini):** Kelimenin Türkçe anlamını yazma ve LLM geri bildirimi.
- [x] **Mod 2 (Cümle Kurma):** Kelimeyi cümlede kullanma, LLM ile gramer ve doğallık analizi.
- [x] **Kütüphane Sistemi:**
    - [x] Dinamik liste tarama ve tikleyerek aktif etme (Multi-select).
    - [x] Yeni liste oluşturma (`CreateListModel`).
    - [x] Global arama (Tüm aktif listelerde tarama ve köken/dosya adı gösterme).
    - [x] Kelime ekleme (Ollama validation + Duplication check + Suggestions).
    - [x] **Merkezi Depolama:** Listeler artık kullanıcı cache dizininde (`~/.cache/foreign7`) tutuluyor.
    - [x] **Export/Dışa Aktar:** Listelerin üzerine gelip 'e' tuşuna basarak dışa aktarma (Yedekleme).
    - [x] **Import/İçe Aktar:** Bilgisayardaki `.txt` dosyalarını kütüphaneye dahil etme.
    - [x] **Liste Silme:** Listeler üzerinden 'x' tuşu ile silme. Son liste silindiğinde yeni liste oluşturmaya zorlar.
- [x] **Teknik Geliştirmeler:**
    - [x] **Gelişmiş Prompting:** `assets/prompts.json` üzerinden sistem talimatlarını dinamik yönetme.
    - [x] **Build Sistemi:** Makefile ile varlıkların (`assets`) otomatik kopyalanması ve `bin/` altına derleme.
    - [x] **Çoklu Dil Desteği:** `assets/lang` altındaki JSON dosyaları ile TR/EN arayüz seçenekleri.
    - [x] **Veri Katmanı Refactoring:** `internal/data` altındaki kodların modüler dosyalara (`config`, `prompts`, `wordlists`) ayrılması.

## 🚀 Sırada Ne Var?

### 1. Yeni Eğitim Modları
- [ ] **Mod 3: Çoktan Seçmeli:** LLM'in ürettiği birbirine yakın 4 şıktan doğrusunu seçme.
- [ ] **Mod 4: Akıllı Sohbet:** Kelime listendeki kelimelere odaklanan, kullanıcıyı onları kullanmaya zorlayan chat modu.

### 2. UX ve Yönetim İyileştirmeleri
- [ ] **Kelime Silme:** Kütüphane içinden münferit kelime silme.
- [ ] **Skor Ekranı:** Test sonunda başarı yüzdesi ve yanlışların özeti.
- [ ] **Aralıklı Tekrar (Spaced Repetition):** Yanlış bilinen kelimeleri daha sık sorma mantığı.

### 3. Teknik Geliştirmeler
- [ ] **Ses Desteği:** Kelimelerin telaffuzlarını dinleme özelliği.

## 🛠 Mevcut Mimari Notları
- **Ana Dosya:** `cmd/foreign7/main.go` (Router tanımları ve başlatma)
- **Veri Katmanı:** `internal/data/*.go` (Config, Prompts, Wordlists ve Store olarak ayrıştırıldı)
- **Veri Yolu:** `~/.cache/foreign7/*.txt` (Doğrudan cache kullanımı)
- **Modeller:** `internal/model/*.go` (Her ekran bir model: menu, test, sentence_mode, manager, list_select, create_list, import_list, settings)
- **Model:** `translategemma:latest` (Ollama)
- **Promptlar:** `assets/prompts.json` (JSON tabanlı template sistemi)
- **Dil Dosyaları:** `assets/lang/*.json` (Arayüz çevirileri)

