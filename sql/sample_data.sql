INSERT INTO scout_group (id, name) VALUES (1, 'Boschetto misterioso');
INSERT INTO scout_group (id, name) VALUES (2, '富士のスカウト');

INSERT INTO account (id, name, email, password, address) VALUES
  (1, 'Giangi', 'giangi@scout.hhh', 'cantlogin', 'Paese'),
  (2, 'Scubi', 'scubi@scout.hhh', 'cantlogin', 'Fori porta'),
  (3, 'Ninni', 'ninni@scout.hhh', 'cantlogin', 'Frazione'),
  (4, 'Lallo', 'lallo@scout.hhh', 'cantlogin', 'Frazione'),
  (5, 'Mimmo', 'mimmo@scout.hhh', 'cantlogin', 'Città'),
  (6, '大翔', 'haruto@scout.zzz', 'cantlogin', '東京'),
  (7, '陽葵', 'himari@scout.zzz', 'cantlogin', '東京');

INSERT INTO account_grant (permission_id, account_id, group_id) VALUES
  (1, 1, 1),
  (1, 2, 1),
  (1, 3, 1),
  (1, 4, 1),
  (1, 5, 1),
  (1, 6, 2),
  (1, 7, 2);

INSERT INTO scout (id, name, group_id) VALUES
  (1, 'Goso', 1),
  (2, 'Gosa', 1),
  (3, 'Sisto', 1),
  (4, 'Noso', 1),
  (5, 'Nosa', 1),
  (6, 'Laso', 1),
  (7, 'Lasa', 1),
  (8, 'Moso', 1),
  (9, 'Mosa', 1),
  (10, 'Mose', 1);

INSERT INTO tutor_scout (tutor_id, scout_id) VALUES
  (1, 1),
  (1, 2),
  (2, 3),
  (3, 4),
  (3, 5),
  (4, 6),
  (4, 7),
  (5, 7),
  (5, 9),
  (5, 10);

INSERT INTO program_activity (id, date, `from`, `to`, location) VALUES
  (1, '2018-12-26', '16:00:00', '23:00:00', 'Piazza'),
  (2, '2019-03-20', '16:00:00', '18:30:00', 'Tana');
