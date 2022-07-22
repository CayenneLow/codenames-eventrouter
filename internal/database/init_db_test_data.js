db.createCollection('test')

db.test.insertOne({
  "GameID": "T35T1",
  "type": "joinGame",
  "sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa2",
  "timestamp": 1658494936,
  "payload": {
    "status": "",
    "message": {
      "clientType": "server"
    }
  }
})

db.test.insertOne({
  "GameID": "T35T2",
  "type": "joinGame",
  "sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa2",
  "timestamp": 1658494936,
  "payload": {
    "status": "",
    "message": {
      "clientType": "server"
    }
  }
})

db.test.insertOne({
  "GameID": "T35T2",
  "type": "joinGame",
  "sessionID": "18c7c74a-317f-46d5-aac8-34a629d82fa3",
  "timestamp": 1658494937,
  "payload": {
    "status": "",
    "message": {
      "clientType": "spymaster"
    }
  }
})