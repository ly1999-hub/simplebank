-- name: CreateAccount :one
INSERT INTO accounts (
 owner,balance,curency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts 
where id=$1 limit 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts 
where id=$1 limit 1
FOR UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE accounts.owner=$1
ORDER BY id
limit $2
OFFSET $3;

-- name: UpdateAccount :one
Update accounts
SET balance=$2
where id=$1
RETURNING *;

-- name: DeleteAccount :exec
Delete From accounts where id=$1;