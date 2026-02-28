# Register

curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice11","password":"secret123"}'

# Login – stores cookies in cookies.txt

curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice1","password":"secret123"}' \
  -c cookies.txt -b cookies.txt

# Profile (authenticated)

curl -X GET http://localhost:8000/profile -b cookies.txt

# Refresh token

curl -X POST http://localhost:8000/refresh -b cookies.txt -c cookies.txt

# Logout

curl -X POST http://localhost:8000/logout -b cookies.txt -c cookies.txt
