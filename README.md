# BeatScorer

## Installation
\`\`\`bash
git clone https://github.com/HocineDoukali/API-DIGIMON/blob/main/README.md
cd beatscorer
go mod download
go run main.go
\`\`\`

## Architecture
- models/ : Structures de données
- services/ : Appels API
- controllers/ : Logique métier
- templates/ : Vues HTML

## Endpoints API utilisés
- GET /player/{id}/full
- GET /player/{id}/scores
- GET /leaderboards?search={query}