SET SESSION default_storage_engine = "MyISAM";
SET SESSION time_zone = "+0:00";
ALTER DATABASE CHARACTER SET "utf8";


CREATE TABLE scout_group (
  id INTEGER NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,

  PRIMARY KEY (id)
);


CREATE TABLE permission (
  id INTEGER NOT NULL,

  PRIMARY KEY (id)
);

INSERT INTO permission (id) VALUES
-- code    label               description                                            table
-- ----    -----               -----------                                            -----
   (1), -- member:             has scouts, can partecipate in excurtion coordination  account_roles
   (2), -- excursion_manager:  can add new excursion                                  account_roles
   (3)  -- group_admin:        can edit group info                                    account_roles
;


CREATE TABLE account (
  id INTEGER NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL DEFAULT '',

  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL DEFAULT '',

  address VARCHAR(255) NOT NULL DEFAULT '',

  PRIMARY KEY (id)
);


CREATE TABLE account_grant (
  permission_id INT NOT NULL,
  account_id INT NOT NULL,
  group_id INT NOT NULL,

  FOREIGN KEY (permission_id) REFERENCES permission(id) ON DELETE CASCADE,
  FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
  FOREIGN KEY (group_id) REFERENCES scout_group(id) ON DELETE CASCADE,
  PRIMARY KEY (permission_id, account_id, group_id)
);


CREATE TABLE invitation (
  token VARCHAR(255) NOT NULL, -- a random uuid
  email VARCHAR(255) NOT NULL, -- the email the invitation has been sent to
  created_on TIMESTAMP NOT NULL,
  -- expiration date is calculated from created_on + SOME SETTINGS

  group_id INT NOT NULL,

  PRIMARY KEY (token)
);


CREATE TABLE scout (
  id INTEGER NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  group_id INT NOT NULL,

  FOREIGN KEY (group_id) REFERENCES scout_group(id) ON DELETE CASCADE,
  PRIMARY KEY (id)
);


CREATE TABLE tutor_scout (
  tutor_id INTEGER NOT NULL,
  scout_id INTEGER NOT NULL,

  FOREIGN KEY (tutor_id) REFERENCES account(id) ON DELETE CASCADE,
  FOREIGN KEY (scout_id) REFERENCES scout(id) ON DELETE CASCADE,
  PRIMARY KEY (tutor_id, scout_id)
);


CREATE TABLE program_activity (
  id INTEGER NOT NULL AUTO_INCREMENT,

  `date` DATE NOT NULL DEFAULT '2010-01-01',
  `from` TIME NOT NULL DEFAULT '00:00:00',
  `to` TIME NOT NULL DEFAULT '00:00:00',
  location VARCHAR(255) NOT NULL DEFAULT '',

  PRIMARY KEY (id)
);


CREATE TABLE scout_activity (
  activity_id INTEGER NOT NULL,
  scout_id INTEGER NOT NULL,
  participate TINYINT NOT NULL DEFAULT 1,

  FOREIGN KEY (activity_id) REFERENCES program_activity(id) ON DELETE CASCADE,
  FOREIGN KEY (scout_id) REFERENCES scout(id) ON DELETE CASCADE
);


CREATE TABLE tutor_activity (
  activity_id INTEGER NOT NULL,
  account_id INTEGER NOT NULL,

  direction CHAR(1) NOT NULL DEFAULT 'O', -- Out, Return
  role CHAR(1) NOT NULL DEFAULT 'N', -- None, Free seats, Rider
  free_seats INTEGER NOT NULL DEFAULT 0,

  FOREIGN KEY (activity_id) REFERENCES program_activity(id) ON DELETE CASCADE,
  FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
  PRIMARY KEY (activity_id, account_id, direction)
);


CREATE TABLE booking (
  id INTEGER NOT NULL AUTO_INCREMENT,

  requester_id INTEGER NOT NULL, -- who needs his/her children to be taken
  rider_id INTEGER NOT NULL, -- who transport children

  state CHAR(1) DEFAULT 'N', -- Not confirmed, Accepted, Revoked / cancelled
  meeting_place VARCHAR(255) NOT NULL DEFAULT '',
  meeting_time TIME NOT NULL,

  PRIMARY KEY (id)
);

CREATE TABLE booking_chat_message (
  id INTEGER NOT NULL AUTO_INCREMENT,

  booking_id INTEGER NOT NULL,
  author_id INTEGER NOT NULL,
  msg VARCHAR(255) NOT NULL,
  msg_date DATETIME NOT NULL,

  PRIMARY KEY (id)
);
