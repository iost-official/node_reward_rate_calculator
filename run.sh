node real.js > real.csv
node theory.js > votes.go
go run ./... | tail -n +2 | tac > theory.csv
zip vote_return_`date +%Y%m%d`.zip *.csv
