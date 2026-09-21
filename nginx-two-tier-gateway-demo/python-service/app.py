import os
from flask import Flask, jsonify, request
app = Flask(__name__)
items = {1: {"id": 1, "name": "Python item"}}
next_id = 2

@app.get('/health')
def health(): return jsonify(status='ok', service='python', instance=os.getenv('INSTANCE_NAME','python-service'))
@app.get('/whoami')
def whoami(): return jsonify(service='python', instance=os.getenv('INSTANCE_NAME','python-service'))
@app.get('/items')
def list_items(): return jsonify(list(items.values()))
@app.post('/items')
def create_item():
    global next_id
    data=request.get_json(silent=True) or {}
    if not data.get('name'): return jsonify(error='name is required'),400
    item={'id':next_id,'name':data['name']}; items[next_id]=item; next_id+=1
    return jsonify(item),201
@app.get('/items/<int:item_id>')
def get_item(item_id):
    item=items.get(item_id)
    return (jsonify(item),200) if item else (jsonify(error='not found'),404)
@app.put('/items/<int:item_id>')
def update_item(item_id):
    if item_id not in items: return jsonify(error='not found'),404
    data=request.get_json(silent=True) or {}
    if not data.get('name'): return jsonify(error='name is required'),400
    items[item_id]['name']=data['name']; return jsonify(items[item_id])
@app.delete('/items/<int:item_id>')
def delete_item(item_id):
    if item_id not in items: return jsonify(error='not found'),404
    del items[item_id]; return '',204
