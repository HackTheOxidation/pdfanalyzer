from PyPDF2 import PdfReader

reader = PdfReader("../../assets/gcc.pdf")

for page in reader.pages:
    text = page.extract_text()
