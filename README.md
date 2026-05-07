# Foreign7 - Language Learning Tool 🌍

[![Go Version](https://img.shields.io/github/go-mod/go-version/ichbinbekir/foreign7)](https://golang.org)
[![Ollama](https://img.shields.io/badge/LLM-Ollama-blue)](https://ollama.ai)

> [English](#english) | [Türkçe](#türkçe)

---

<a name="english"></a>
## 🇬🇧 English

Foreign7 is a Terminal User Interface (TUI) application designed to help you master new languages using the power of Local LLMs (via Ollama). It functions as a specialized tool that validates your vocabulary, tests your meaning predictions, and evaluates your sentence-building skills.

### ✨ Key Features

*   **Dual Learning Modes:**
    *   **Meaning Prediction:** Guess the meaning of words in your target language and get instant AI feedback.
    *   **Sentence Building:** Practice using words in context. The AI analyzes your grammar, meaning, and provides natural alternatives.
*   **Smart Library Management:**
    *   Organize words into multiple categories (lists).
    *   Search across all active lists simultaneously.
    *   Import and Export your word lists as `.txt` files.
*   **Privacy First & Local AI:** Powered by Ollama. Your data and learning history stay on your machine.
*   **Centralized Storage:** Automatically stores your data in the system's cache directory (`~/.cache/foreign7`).
*   **Multilingual UI:** Supports both English and Turkish interfaces.
*   **Modern TUI:** Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for a sleek, responsive terminal experience.

### 🛠 Prerequisites

1.  **Go:** Version 1.21 or higher.
2.  **Ollama:** Installed and running.
3.  **Model:** Pull the required model:
    ```bash
    ollama pull translategemma:latest
    ```

### 🚀 Installation & Build

1.  Clone the repository:
    ```bash
    git clone https://github.com/ichbinbekir/foreign7.git
    cd foreign7
    ```
2.  Build the project using the Makefile:
    ```bash
    make build
    ```
    This will create a `bin/` directory containing the executable and required assets.

3.  Run the application:
    ```bash
    ./bin/foreign7
    ```

### ⌨️ Key Bindings

*   **Enter:** Confirm / Select / Next
*   **Esc:** Go back / Cancel
*   **Space:** Toggle list active status (in Library)
*   **'e':** Export list (in Library)
*   **'x':** Delete list (in Library)
*   **Ctrl+C:** Quit

---

<a name="türkçe"></a>
## 🇹🇷 Türkçe

Foreign7, Yerel YZ (Ollama) gücünü kullanarak yeni diller öğrenmenize yardımcı olan bir Terminal Kullanıcı Arayüzü (TUI) uygulamasıdır. Kelime bilginizi doğrulayan, anlam tahminlerinizi test eden ve cümle kurma becerilerinizi değerlendiren uzman bir dil öğrenme aracı gibi çalışır.

### ✨ Temel Özellikler

*   **İki Farklı Eğitim Modu:**
    *   **Anlam Tahmini:** Hedef dildeki kelimelerin anlamını tahmin edin ve anında YZ geri bildirimi alın.
    *   **Cümle Kurma:** Kelimeleri bağlam içinde kullanma pratiği yapın. YZ gramerinizi ve anlamınızı analiz eder, daha doğal öneriler sunar.
*   **Akıllı Kütüphane Yönetimi:**
    *   Kelimeleri kategorilere (listelere) ayırın.
    *   Aktif olan tüm listelerde aynı anda arama yapın.
    *   Kelime listelerinizi `.txt` dosyası olarak İçe/Dışa Aktarın.
*   **Gizlilik Odaklı & Yerel YZ:** Ollama ile çalışır. Verileriniz ve öğrenme geçmişiniz bilgisayarınızda kalır.
*   **Merkezi Depolama:** Verilerinizi otomatik olarak sistemin cache dizininde (`~/.cache/foreign7`) tutar.
*   **Çok Dilli Arayüz:** Hem İngilizce hem de Türkçe arayüz desteği mevcuttur.
*   **Modern TUI:** Şık ve hızlı bir terminal deneyimi için [Bubble Tea](https://github.com/charmbracelet/bubbletea) ile geliştirilmiştir.

### 🛠 Gereksinimler

1.  **Go:** Versiyon 1.21 veya üzeri.
2.  **Ollama:** Bilgisayarınızda kurulu ve çalışır durumda olmalı.
3.  **Model:** Gerekli modeli indirin:
    ```bash
    ollama pull translategemma:latest
    ```

### 🚀 Kurulum ve Derleme

1.  Depoyu klonlayın:
    ```bash
    git clone https://github.com/ichbinbekir/foreign7.git
    cd foreign7
    ```
2.  Makefile kullanarak derleyin:
    ```bash
    make build
    ```
    Bu komut, çalıştırılabilir dosyayı ve gerekli varlıkları içeren bir `bin/` dizini oluşturacaktır.

3.  Uygulamayı çalıştırın:
    ```bash
    ./bin/foreign7
    ```

### ⌨️ Tuş Kombinasyonları

*   **Enter:** Onayla / Seç / Sonraki
*   **Esc:** Geri dön / Vazgeç
*   **Boşluk (Space):** Listeyi aktif/pasif yap (Kütüphane ekranında)
*   **'e':** Listeyi dışa aktar (Kütüphane ekranında)
*   **'x':** Listeyi sil (Kütüphane ekranında)
*   **Ctrl+C:** Çıkış

---

## 📜 License / Lisans

Distributed under the MIT License. See `LICENSE` for more information.
MIT Lisansı ile dağıtılmaktadır. Daha fazla bilgi için `LICENSE` dosyasına göz atın.
