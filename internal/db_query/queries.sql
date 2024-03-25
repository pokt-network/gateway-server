-- name: GetPoktApplications :many
SELECT id, pgp_sym_decrypt(encrypted_private_key, pggen.arg('encryption_key')) AS decrypted_private_key
FROM pokt_applications;

-- name: InsertPoktApplications :exec
INSERT INTO pokt_applications (encrypted_private_key)
VALUES (pgp_sym_encrypt(pggen.arg('private_key'), pggen.arg('encryption_key')));

-- name: DeletePoktApplication :exec
DELETE FROM pokt_applications
WHERE id = pggen.arg('application_id');

-- name: GetChainConfigurations :many
SELECT * FROM chain_configurations;