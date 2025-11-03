from concurrent import futures
import grpc
import job_pb2
import job_pb2_grpc
import logging
from services.resume_parser import ResumeParser
from services.resume_proccessing import ResumeProcessor

MAX_TEXT_PROCESSING_BYTES = 5 * 1024 * 1024

logging.basicConfig(level=logging.INFO)

class JobServiceServicer(job_pb2_grpc.JobServiceServicer):
    def __init__(self):
        logging.info("Initializing JobServiceServicer...")
        self.parser = ResumeParser()
        self.processor = ResumeProcessor(maxi_allowed_bytes=MAX_TEXT_PROCESSING_BYTES)


    def ParseResume(self, request, context):
        # Parsing the resume
        resume_text = self.processor.process_raw_resume(request.resume_file_content,context)
        parsed_data = self.parser.parse(resume_text,context)
        response = job_pb2.ParseResumeResponse()
        response.job_titles = parsed_data.job_titles
        response.skills = parsed_data.skills
        response.experience = parsed_data.experience
        return response

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    job_pb2_grpc.add_JobServiceServicer_to_server(JobServiceServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started, listening on port 50051.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
