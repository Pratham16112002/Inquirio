
import os
import logging
from typing import List
import grpc
from pydantic import BaseModel, Field
from transformers import pipeline, AutoTokenizer, AutoModelForQuestionAnswering
from utils.text_cleaner import clean_text
import re

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# --- Pydantic Model for Structured Output ---
class ResumeData(BaseModel):
    """Structured data extracted from a resume."""
    job_titles: List[str] = Field(default_=[], description="Extracted job titles.")
    skills: List[str] = Field(default_=[], description="Extracted skills.")
    experience: str = Field(default="", description="A summary of the work experience.")

# --- Resume Parser Service ---
class ResumeParser:
    """
    A service to parse resumes and extract key information using a Hugging Face
    Question-Answering model.
    """
    def __init__(self, model_name: str = None):
        """
        Initializes the parser and loads the QA model.
        
        Args:
            model_name (str, optional): The name of the Hugging Face model to use.
                                        Defaults to the one specified in the
                                        RESUME_MODEL environment variable or a default.
        """
        try:
            model_name = model_name or os.getenv("RESUME_MODEL", "distilbert-base-cased-distilled-squad")
            logger.info(f"Loading model: {model_name}")

            # Load the model and tokenizer
            tokenizer = AutoTokenizer.from_pretrained(model_name)
            model = AutoModelForQuestionAnswering.from_pretrained(model_name)

            # Set up the QA pipeline
            self.nlp = pipeline("question-answering", model=model, tokenizer=tokenizer)
            logger.info("ResumeParser initialized successfully.")

        except Exception as e:
            logger.error(f"Failed to load model or tokenizer: {e}", exc_info=True)
            raise

    def _query_model(self, question: str, context: str) -> str:
        """Sends a question and context to the QA model and returns the answer."""
        try:
            result = self.nlp(question=question, context=context)
            return result['answer']
        except Exception as e:
            logger.error(f"Error during model query for question '{question}': {e}", exc_info=True)
            return ""

    def parse(self, resume_text: str,context) -> ResumeData:
        """
        Parses the raw text of a resume to extract structured data.

        Args:
            resume_text (str): The full text content of the resume.

        Returns:
            ResumeData: A Pydantic model containing the extracted information.
        """
        if not resume_text:
            logger.warning("Input resume text is empty.")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("File has nothing to process")
            return None

        # 1. Clean the text using the utility
        cleaned_context = clean_text(resume_text)
        
        logger.info("Extracting information from resume...")

        # 2. Define questions for the model
        questions = {
            "job_titles": "What are the job titles?",
            "skills": "What are the skills?",
            "experience": "Summarize the work experience."
        }

        # 3. Query the model for each piece of information
        extracted_data = {}
        for key, question in questions.items():
            answer = self._query_model(question, cleaned_context)
            
            # Simple post-processing for list-based fields
            if key in ["job_titles", "skills"] and answer:
                # Split by common delimiters and clean up
                items = re.split(r',|\n|;', answer)
                extracted_data[key] = [item.strip() for item in items if item.strip()]
            else:
                extracted_data[key] = answer

        logger.info(f"Successfully extracted: {extracted_data}")

        # 4. Return structured data
        return ResumeData(**extracted_data)

# --- Example Usage ---
if __name__ == '__main__':
    # This block allows for direct testing of the parser
    
    # Ensure the utility is available in the path
    import sys
    sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
    from utils.text_cleaner import clean_text
    import re

    # Example resume text (replace with a real one for better testing)
    sample_resume = """
    John Doe
    Software Engineer

    Experience:
    - Senior Developer at Tech Corp (2020-Present)
      - Led a team to build a new microservice using Python, Django, and Kubernetes.
    - Software Engineer at Innovate LLC (2018-2020)
      - Developed features for a web application with React and Node.js.

    Skills:
    Python, Java, C++, JavaScript, React, Node.js, Docker, Kubernetes, SQL, Git.
    """

    print("Initializing ResumeParser...")
    parser = ResumeParser()
    
    print("\n--- Parsing Sample Resume ---")
    parsed_data = parser.parse(sample_resume)

    print("\n--- Extracted Data ---")
    print(f"Job Titles: {parsed_data.job_titles}")
    print(f"Skills: {parsed_data.skills}")
    print(f"Experience Summary: {parsed_data.experience}")
    print("\n----------------------")

