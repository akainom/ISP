curl -X POST http://10.186.18.43:8080/celebrities -H "Content-Type: application/json" -d '{"fullName": "Паша Техник", "nationality": "RU", "reqPhotoPath": "./pablito.jpg"}'
curl -X POST http://10.186.18.43:8080/celebrities -H "Content-Type: application/json" -d '{"fullName": "1.kla$", "nationality": "DE", "reqPhotoPath": "./temy4.jpg"}'
curl -X POST http://10.186.18.43:8080/celebrities -H "Content-Type: application/json" -d '{"fullName": "miron daun", "nationality": "RU", "reqPhotoPath": "./miron.jpg"}'
curl -X PUT http://10.186.18.43:8080/celebrity/2 -H "Content-Type: application/json" -d '{"id": 1, "fullName": "Павел Техник", "nationality": "RU", "reqPhotoPath": "/img/new_pasha.jpg"}'
curl -X GET http://10.186.18.43:8080/celebrities

