node 实际.js > 实际.csv
node 理论.js > votes.go
go run ./... | tail -n +2 | tac > 理论.csv
zip vote_return_`date +%Y%m%d`.zip *.csv
