for i in {1..10}; do 
  echo "Attempt $i"
  curl -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"bad"}' http://localhost:3000/login
done
