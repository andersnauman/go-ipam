DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id            INT AUTO_INCREMENT NOT NULL,
  username      VARCHAR(64) NOT NULL,
  password_hash VARCHAR(128) NOT NULL,
  password_salt VARCHAR(16) NOT NULL,
  first_name    VARCHAR(32) NOT NULL,
  last_name     VARCHAR(32) NOT NULL,
  PRIMARY KEY ('id')   
);

INSERT INTO users
  (username, password_hash, password_salt, first_name, last_name)
VALUES
  ('admin', 'b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86', '123456', 'Admin', 'Adminsson');

DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
  id          INT AUTO_INCREMENT NOT NULL,
  user_id     INT NOT NULL,
  cookie_hash VARCHAR(128) NOT NULL,
  cookie_salt VARCHAR(64) NOT NULL,
  created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
  expire      DATETIME NOT NULL,
  PRIMARY KEY ('id')   
);

DROP TABLE IF EXISTS subnets;
CREATE TABLE subnets (
  subnet      VARCHAR(128) NOT NULL,
  description VARCHAR(256),
  vrf_id      INT,
  vlan_id     INT,
  site_id     INT,
  city        VARCHAR(64),
  country     VARCHAR(64),
  coordinates VARCHAR(32),
  PRIMARY KEY ('subnet')
);

INSERT INTO subnets
  (subnet, description, vrf_id, vlan_id, site_id, city, country, coordinates)
VALUES
  ('192.168.0.0/24', "Mgmt-1", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.1.0/24', "Mgmt-2", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.2.0/24', "Mgmt-3", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.3.0/24', "Mgmt-4", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.4.0/24', "Mgmt-5", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.5.0/24', "Mgmt-6", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.6.0/24', "Mgmt-7", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.7.0/24', "Mgmt-8", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240"),
  ('192.168.8.0/24', "Mgmt-9", 1, 1, 1, "stockholm", "sweden", "59.334591, 18.063240");

DROP TABLE IF EXISTS devices;
CREATE TABLE devices (
  ip_address  VARCHAR(64) NOT NULL,
  description VARCHAR(256),
  hostname    VARCHAR(64),
  status      INT NOT NULL,
  subnet_id   INT NOT NULL,
  PRIMARY KEY ('ip_address')
);

INSERT INTO devices
  (ip_address, description, hostname, status, subnet_id)
VALUES
  ('192.168.0.1', 'dns-server', 'host1.example.com', 1, 1),
  ('192.168.0.2', 'dns-server', 'host2.example.com', 1, 1),
  ('192.168.0.3', 'dns-server', 'host3.example.com', 1, 1),
  ('192.168.0.4', 'dns-server', 'host4.example.com', 1, 1);