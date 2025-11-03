import grpc
import logging
import io
from pdfminer.high_level import extract_text as extract_pdf_text


# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class ResumeProcessor:
    def __init__(self,maxi_allowed_bytes):
        self.maxi_allowed_bytes = maxi_allowed_bytes
    

    def process_raw_resume(self,file_bytes : bytes, context ) -> str:

        if len(file_bytes) > self.maxi_allowed_bytes:
            logging.info("File exceeded the maximum allowed bytes")
            context.abort(grpc.StatusCode.RESOURCE_EXHAUSTED, "File too large")
            return
        
        # Convertion of the raw bytes to a string
        try:
            """
            Wrapping the raw bytes in a BytesIO object , This creates
            an in memory binary file that the pdfminer can read
            """
            pdf_file_in_memory = io.BytesIO(file_bytes)

            """
            Pass the in memory binary file to the pdfminer extractor
            """
            text = extract_pdf_text(pdf_file_in_memory)
        except Exception as e:
            logging.info("Error failed to parse the file content {e}")
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "File not in UTF-8 format")
            return
        

        return text

