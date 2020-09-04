"""__init__.py is a constructor document that lets the python interpreter 
    see the folder named app as if it were a package that we can import objects from

    First import the Flask class from the flask library
    __name__ refers to the name in reference to the main python process

    Next import the pages of the web app stored in routes
    """
from app import routes
from flask import Flask

app = Flask(__name__)
