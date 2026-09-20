from flask import Flask, request, jsonify
import socket, os
app = Flask(__name__)

@app.get('/')
def home():
    return jsonify(service='python-service', hostname=socket.gethostname(), path=request.path, mode=os.getenv('MODE','demo'))

@app.get('/python')
def python():
    return jsonify(service='python-service', hostname=socket.gethostname(), path=request.path)

@app.get('/health')
def health(): return 'OK', 200

app.run(host='0.0.0.0', port=8080)
