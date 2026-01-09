CREATE TABLE cpus
(
    cpu_id                     BIGSERIAL PRIMARY KEY,
    name                       VARCHAR(128) NOT NULL,
    brand                      VARCHAR(64),
    socket_code                VARCHAR(32)  NOT NULL REFERENCES cpu_sockets (socket_code),
    architecture               VARCHAR(64),
    core_count                 INTEGER,
    thread_count               INTEGER,
    base_clock_mhz             INTEGER,
    boost_clock_mhz            INTEGER,
    tdp_watt                   INTEGER,
    has_integrated_gpu         BOOLEAN,
    supported_ram_type         VARCHAR(16), -- DDR4 / DDR5
    supported_ram_freq_max_mhz INTEGER,
    memory_channels            INTEGER,
    pcie_version               INTEGER,
    pcie_lanes_total           INTEGER
);