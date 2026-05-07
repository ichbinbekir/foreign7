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

## 🚀 Sırada Ne Var?

### 1. Yeni Eğitim Modları
- [ ] **Mod 3: Çoktan Seçmeli:** LLM'in ürettiği birbirine yakın 4 şıktan doğrusunu seçme.
- [ ] **Mod 4: Akıllı Sohbet:** Kelime listendeki kelimelere odaklanan, kullanıcıyı onları kullanmaya zorlayan chat modu.

### 2. UX ve Yönetim İyileştirmeleri
- [ ] **Kelime Silme:** Kütüphane içinden münferit kelime silme.
- [ ] **Skor Ekranı:** Test sonunda başarı yüzdesi ve yanlışların özeti.
- [ ] **Aralıklı Tekrar (Spaced Repetition):** Yanlış bilinen kelimeleri daha sık sorma mantığı.

### 3. Teknik Geliştirmeler
- [ ] **Gelişmiş Prompting:** `prompts.json` üzerinden sistem talimatlarını daha dinamik yönetme.
- [ ] **Ses Desteği:** Kelimelerin telaffuzlarını dinleme özelliği.

## 🛠 Mevcut Mimari Notları
- **Ana Dosya:** `cmd/foreign7/main.go` (Router tanımları ve başlatma)
- **Veri Yolu:** `~/.cache/foreign7/*.txt` (Doğrudan cache kullanımı)
- **Modeller:** `internal/model/*.go` (Her ekran bir model: menu, test, sentence_mode, manager, list_select, create_list, import_list)
- **Model:** `translategemma:latest` (Ollama)
