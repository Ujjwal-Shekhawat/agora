source ./.env
docker exec -it vault vault operator unseal ${UNSEAL_KEY}