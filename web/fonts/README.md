# Local Fonts

This app prefers locally hosted fonts for offline use and strict CSPs.

Expected files (place in this directory):

- Inter-Variable.woff2 — Latin UI text (variable weight 100–900)
- Amiri-Regular.woff2 — Arabic body
- Amiri-Bold.woff2 — Arabic headings / emphasis
- NotoNaskhArabic-Regular.woff2 — Arabic body (alternate)
- NotoNaskhArabic-SemiBold.woff2 — Arabic emphasis (alternate)

Recommended sources and licenses:

- Inter: https://github.com/rsms/inter (OFL-1.1). Release assets typically include `Inter-Variable.woff2`.
- Amiri: https://github.com/aliftype/amiri (OFL-1.1). Release assets include `Amiri-Regular.woff2` and `Amiri-Bold.woff2`.
- Noto Naskh Arabic: https://github.com/notofonts/arabic (OFL-1.1). Look for packaged WOFF2 files for Noto Naskh Arabic Regular and SemiBold.

After placing files, no further config is needed — `web/styles.css` declares `@font-face` rules pointing to these filenames.

Notes
- If you prefer Noto Naskh Arabic, you can add `NotoNaskhArabic-Regular.woff2` and `NotoNaskhArabic-SemiBold.woff2` and update `styles.css` to reference them.
- For best performance, keep to WOFF2 format and subset fonts if you target specific Unicode ranges.
