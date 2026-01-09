CREATE TABLE cpu_sockets
(
    socket_code VARCHAR(32) PRIMARY KEY,
    description TEXT
);

CREATE TABLE motherboard_form_factors
(
    form_factor_code VARCHAR(16) PRIMARY KEY, -- ATX, mATX, ITX
    description      TEXT,
    width_mm         INTEGER,
    height_mm        INTEGER
);