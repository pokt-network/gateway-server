-- Drop the table 'pokt_applications' and its dependencies
DROP TABLE IF EXISTS pokt_applications;

-- Drop the table 'base_model'
DROP TABLE IF EXISTS base_model;

-- Drop the extensions 'pgcrypto' and 'uuid-ossp'
DROP EXTENSION IF EXISTS pgcrypto;
DROP EXTENSION IF EXISTS "uuid-ossp";
