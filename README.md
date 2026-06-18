# go-ocr

Công cụ CLI viết bằng Go để chuyển đổi tài liệu **PDF / Word (.docx)** sang **Markdown (.md)**.

Hỗ trợ:
- Tự động phát hiện: tài liệu có sẵn text thì trích xuất trực tiếp; tài liệu scan/ảnh thì OCR.
- OCR theo hướng hybrid: Tesseract trích xuất trước, dùng LLM (chuẩn OpenAI-compatible) để tinh chỉnh thành Markdown và làm fallback.
- Mô tả ảnh nhúng trong tài liệu bằng LLM (nếu model hỗ trợ vision).

## Yêu cầu môi trường

- [Go](https://go.dev/dl/) 1.22+ (đang dùng 1.26.x)
- [Tesseract OCR](https://github.com/UB-Mannheim/tesseract/wiki) (kèm `vie.traineddata` nếu cần OCR tiếng Việt)
- [Poppler](https://github.com/oschwartz10612/poppler-windows/releases) (cung cấp `pdftoppm` để chuyển PDF sang ảnh)
- (Tùy chọn) API LLM chuẩn OpenAI-compatible: base URL + API key + tên model

Tất cả binary ngoài (`tesseract`, `pdftoppm`) cần nằm trong PATH.

## Cài đặt

```bash
git clone git@github.com:lbkkk/OCR_go.git
cd OCR_go
go build ./...
```

## Sử dụng

```bash
goocr convert <đường-dẫn-file-hoặc-thư-mục> [flags]
```

Các flag dự kiến:
- `--out`        : thư mục/đường dẫn xuất file `.md`
- `--lang`       : ngôn ngữ OCR (mặc định `eng+vie`)
- `--force-ocr`  : ép OCR kể cả khi có text layer
- `--use-llm`    : bật tinh chỉnh + mô tả ảnh bằng LLM
- `--recursive`  : xử lý toàn bộ thư mục

## Cấu trúc dự án

Theo chuẩn [golang-standards/project-layout](https://github.com/golang-standards/project-layout). Xem chi tiết trong `docs/`.

## Giấy phép

Chưa xác định.
