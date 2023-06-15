create table if not exists users
(
    userid SERIAL primary key,
    login    varchar(100),
    password varchar(100),
    role     varchar(100)
);

create table if not exists cooling
(
    coolingid         SERIAL primary key,
    name              varchar(100),
    brand             varchar(100),
    max_speed         int,
    count_ventilators int,
    price             int
);

create table if not exists cpu
(
    cpuid         SERIAL primary key,
    name          varchar(100),
    brand         varchar(100),
    series        varchar(100),
    model    varchar(100),
    basic_frequency float,
    supported_ddr varchar(100),
    min_frequency_of_ram int,
    max_frequency_of_ram int,
    price         int
);

create table if not exists motherboard
(
    motherboardid SERIAL primary key,
    name          varchar(100),
    brand         varchar(100),
    count_slots   int,
    supported_ddr varchar(100),
    form_factor_of_ram varchar(100),
    min_frequency_of_ram int,
    max_frequency_of_ram int,
    price         int
);


create table if not exists pc_build
(
    pc_buildid SERIAL primary key,
    user_id INT,
    cpu_id INT,
    ram_id INT,
    powerbox_id int,
    motherboard_id int,
    cooling_id int,
    videocard_id int,
    compatibility varchar(100),
    price int,
    FOREIGN KEY (user_id) REFERENCES users(userid),
    FOREIGN KEY (cpu_id) REFERENCES cpu(cpuid),
    FOREIGN KEY (ram_id) REFERENCES ram(ramid),
    FOREIGN KEY (powerbox_id) REFERENCES powerbox(powerboxid),
    FOREIGN KEY (motherboard_id) REFERENCES motherboard(motherboardid),
    FOREIGN KEY (cooling_id) REFERENCES cooling(coolingid),
    FOREIGN KEY (videocard_id) REFERENCES videocard(videocardid)
);

create table if not exists powerbox
(
    powerboxid  SERIAL primary key,
    name        varchar(100),
    brand       varchar(100),
    power       int,
    form_factor varchar(100),
    price       int
);

create table if not exists ram
(
    ramid     SERIAL primary key,
    name      varchar(100),
    brand     varchar(100),
    rammemory int,
    ddr       varchar(100),
    form_factor varchar(100),
    clock_frequency int,
    price     int
);

create table if not exists videocard
(
    videocardid SERIAL primary key,
    name        varchar(100),
    brand       varchar(100),
    series      int,
    vmemory     int,
    price       int
);

CREATE TABLE IF NOT EXISTS cpu_ram
(
    cpu_id SERIAL REFERENCES cpu (cpuid),
    ram_id SERIAL REFERENCES ram (ramid),
    PRIMARY KEY (cpu_id, ram_id)
);

CREATE TABLE IF NOT EXISTS cpu_motherboard
(
    cpu_id SERIAL REFERENCES cpu (cpuid),
    motherboard_id SERIAL REFERENCES motherboard (motherboardid),
    PRIMARY KEY (cpu_id, motherboard_id)
);

CREATE TABLE IF NOT EXISTS motherboard_ram
(
    motherboard_id SERIAL REFERENCES motherboard (motherboardid),
    ram_id SERIAL REFERENCES ram (ramid),
    PRIMARY KEY (motherboard_id, ram_id)
);

-- create view cpu_view as
-- SELECT name, price
-- FROM cpu;
--
-- create view cooling_view as
-- SELECT name, price
-- FROM cpu;
--
-- create view motherboard_view as
-- SELECT name, price
-- FROM motherboard;
--
-- create view powerbox_view as
-- SELECT name, price
-- FROM powerbox;
--
-- create view ram_view as
-- SELECT name, price
-- FROM ram;
--
-- create view videocard_view as
-- SELECT name, price
-- FROM videocard;

--------------------------------------------------создание ролей
CREATE ROLE "guest";
CREATE ROLE "user";
CREATE ROLE "admin";--роли уже созданы
--------------------------------------------------

--------------------------------------------------наделение админа правами суперюзера
ALTER ROLE "admin" WITH SUPERUSER;
--------------------------------------------------

--------------------------------------------------права для гестов
GRANT SELECT ON TABLE cpu TO "guest";
GRANT SELECT ON TABLE cooling TO "guest";
GRANT SELECT ON TABLE motherboard TO "guest";
GRANT SELECT ON TABLE powerbox TO "guest";
GRANT SELECT ON TABLE ram TO "guest";
GRANT SELECT ON TABLE videocard TO "guest";
GRANT INSERT, SELECT ON TABLE users TO "guest";
----------------------------------------------------------------

--------------------------------------------------права для юзеров
GRANT SELECT ON TABLE cpu TO "user";
GRANT SELECT ON TABLE cooling TO "user";
GRANT SELECT ON TABLE motherboard TO "user";
GRANT SELECT ON TABLE powerbox TO "user";
GRANT SELECT ON TABLE ram TO "user";
GRANT SELECT ON TABLE videocard TO "user";
GRANT SELECT ON TABLE cpu_ram TO "user";
GRANT SELECT ON TABLE cpu_motherboard TO "user";
GRANT SELECT ON TABLE motherboard_ram TO "user";
GRANT INSERT, SELECT ON TABLE users TO "user";
GRANT SELECT, UPDATE, INSERT ON TABLE pc_build TO "user";
----------------------------------------------------------------

SELECT * FROM pg_user;
--------------------------------------------------удаление всех прав для user
REVOKE ALL ON TABLE cpu FROM "user";
REVOKE ALL ON TABLE cooling FROM "user";
REVOKE ALL ON TABLE motherboard FROM "user";
REVOKE ALL ON TABLE powerbox FROM "user";
REVOKE ALL ON TABLE ram FROM "user";
REVOKE ALL ON TABLE videocard FROM "user";
REVOKE ALL ON TABLE cpu_ram FROM "user";
REVOKE ALL ON TABLE cpu_motherboard FROM "user";
REVOKE ALL ON TABLE motherboard_ram FROM "user";
REVOKE ALL ON TABLE users FROM "user";
REVOKE ALL ON TABLE pc_build FROM "user";
--------------------------------------------------


--------------------------------------------------удаление всех прав для guest
REVOKE ALL ON TABLE cpu FROM "guest";
REVOKE ALL ON TABLE cooling FROM "guest";
REVOKE ALL ON TABLE motherboard FROM "guest";
REVOKE ALL ON TABLE powerbox FROM "guest";
REVOKE ALL ON TABLE ram FROM "guest";
REVOKE ALL ON TABLE videocard FROM "guest";
REVOKE ALL ON TABLE users FROM "guest";
--------------------------------------------------


--------------------------------------------------удаление ролей
DROP ROLE "guest";
DROP ROLE "user";
DROP ROLE "admin";
--------------------------------------------------


--------------------------------------------------триггер для проверки на совместимость (если частота памяти больше частот проца и материнки то совместимость ок) но в связных таблицах лишь те которые строго в промежутки попадают между мин и макс
CREATE OR REPLACE TRIGGER check_compatibility
    BEFORE INSERT OR UPDATE ON pc_build
    FOR EACH ROW
EXECUTE FUNCTION checkCompatibility();

CREATE OR REPLACE FUNCTION checkCompatibility()
    RETURNS TRIGGER
AS $$
DECLARE
    ram_ddr varchar(100);
    sup_ddr_cpu varchar(100);
    sup_ddr_motherboard varchar(100);
    form_factor_ram varchar(100);
    form_factor_motherboard varchar(100);
    max_frequency_cpu int;
    max_frequency_motherboard int;
    min_frequency_cpu int;
    min_frequency_motherboard int;
    ram_clock_frequency int;
BEGIN
    IF (TG_OP = 'INSERT') OR (TG_OP = 'UPDATE') THEN
        SELECT "ddr" into ram_ddr
        FROM ram
        WHERE ram."ramid" = NEW.ram_id;
        SELECT "supported_ddr" into sup_ddr_cpu
        FROM cpu
        WHERE cpu."cpuid" = NEW.cpu_id;
        SELECT "supported_ddr" into sup_ddr_motherboard
        FROM motherboard
        WHERE motherboard."motherboardid" = NEW.motherboard_id;

        SELECT "form_factor" into form_factor_ram
        FROM ram
        WHERE ram."ramid" = NEW.ram_id;
        SELECT "form_factor_of_ram" into form_factor_motherboard
        FROM motherboard
        WHERE motherboard."motherboardid" = NEW.motherboard_id;

        SELECT "max_frequency_of_ram" into max_frequency_cpu
        FROM cpu
        WHERE cpu."cpuid" = NEW.cpu_id;
        SELECT "max_frequency_of_ram" into max_frequency_motherboard
        FROM motherboard
        WHERE motherboard."motherboardid" = NEW.motherboard_id;

        SELECT "min_frequency_of_ram" into min_frequency_cpu
        FROM cpu
        WHERE cpu."cpuid" = NEW.cpu_id;
        SELECT "min_frequency_of_ram" into min_frequency_motherboard
        FROM motherboard
        WHERE motherboard."motherboardid" = NEW.motherboard_id;
        SELECT "clock_frequency" into ram_clock_frequency
        FROM ram
        WHERE ram."ramid" = NEW.ram_id;

        IF (position(ram_ddr in sup_ddr_cpu) > 0) AND (position(ram_ddr in sup_ddr_motherboard) > 0) AND (form_factor_ram = form_factor_motherboard) AND (ram_clock_frequency >= min_frequency_cpu) AND (ram_clock_frequency >= min_frequency_motherboard)  AND (min_frequency_motherboard <= max_frequency_cpu) AND (min_frequency_cpu <= max_frequency_motherboard) THEN
            NEW.compatibility = 'OK';
        ELSEIF (position(ram_ddr in sup_ddr_cpu) <= 0) THEN
            NEW.compatibility = 'Ram_ddr != sup_ddr_cpu';
        ELSEIF (position(ram_ddr in sup_ddr_motherboard) <= 0) THEN
            NEW.compatibility = 'Ram_ddr != sup_ddr_motherboard';
        ELSEIF (form_factor_ram != form_factor_motherboard) THEN
            NEW.compatibility = 'form_factor_ram != form_factor_motherboard';
        ELSEIF (ram_clock_frequency < min_frequency_cpu) THEN
            NEW.compatibility = 'ram_clock_frequency < min_frequency_cpu';
        ELSEIF (ram_clock_frequency < min_frequency_motherboard) THEN
            NEW.compatibility = 'ram_clock_frequency < min_frequency_motherboard';
        ELSEIF (min_frequency_motherboard > max_frequency_cpu) THEN
            NEW.compatibility = 'min_ram_frequency_motherboard > max_ram_frequency_cpu';
        ELSEIF (min_frequency_cpu > max_frequency_motherboard) THEN
            NEW.compatibility = 'min_ram_frequency_ram > max_ram_frequency_motherboard';
        end if;
    END if;
    RETURN
        NEW;
END;
$$ LANGUAGE plpgsql;
--------------------------------------------------


--------------------------------------------------триггер для подсчета итоговой стоимости сборки
CREATE OR REPLACE TRIGGER sum_price
    BEFORE INSERT OR UPDATE ON pc_build
    FOR EACH ROW
EXECUTE FUNCTION SumPrice();

CREATE OR REPLACE FUNCTION SumPrice()
    RETURNS TRIGGER
AS $$
DECLARE
    ram_price int;
    videocard_price int;
    cooling_price int;
    cpu_price int;
    motherboard_price int;
    powerbox_price int;
BEGIN
    IF (TG_OP = 'INSERT') OR (TG_OP = 'UPDATE') THEN
        SELECT "price" into ram_price
        FROM ram
        WHERE ram."ramid" = NEW.ram_id;
        SELECT "price" into videocard_price
        FROM videocard
        WHERE videocard."videocardid" = NEW.videocard_id;
        SELECT "price" into cooling_price
        FROM cooling
        WHERE cooling."coolingid" = NEW.cooling_id;
        SELECT "price" into cpu_price
        FROM cpu
        WHERE cpu."cpuid" = NEW.cpu_id;
        SELECT "price" into motherboard_price
        FROM motherboard
        WHERE motherboard."motherboardid" = NEW.motherboard_id;
        SELECT "price" into powerbox_price
        FROM powerbox
        WHERE powerbox."powerboxid" = NEW.powerbox_id;
        NEW.price =  ram_price + videocard_price + cooling_price + cpu_price + motherboard_price + powerbox_price;
    end if;
    RETURN
        NEW;
END;
$$ LANGUAGE plpgsql;
--------------------------------------------------


-- SELECT setval('cpu_cpuid_seq', (SELECT MAX(cpuid) FROM cpu));
-- SELECT setval('cooling_coolingid_seq', (SELECT MAX(coolingid) FROM cooling));
-- SELECT setval('powerbox_powerboxid_seq', (SELECT MAX(powerboxid) FROM powerbox));
-- SELECT setval('motherboard_motherboardid_seq', (SELECT MAX(motherboardid) FROM motherboard));
-- SELECT setval('ram_ramid_seq', (SELECT MAX(ramid) FROM ram));
-- SELECT setval('videocard_videocardid_seq', (SELECT MAX(videocardid) FROM videocard));





insert into pc_build(user_id, cpu_id, ram_id, powerbox_id, motherboard_id, cooling_id, videocard_id) VALUES (2, 1, 54, 1, 1, 1, 2);
--------------------------------------------------дроп таблиц
drop table if exists cpu cascade;
drop table if exists cooling cascade;
drop table if exists motherboard cascade;
drop table if exists powerbox cascade;
drop table if exists ram cascade;
drop table if exists videocard cascade;
drop table if exists pc_build cascade;
drop table if exists users cascade;
drop table if exists cpu_ram cascade;
--------------------------------------------------
CREATE EXTENSION pg_trgm;

insert into cpu_ram(cpu_id, ram_id) SELECT cpu.cpuid, ram.ramid FROM cpu, ram WHERE position(ram.ddr in cpu.supported_ddr) > 0 AND ram.clock_frequency >= cpu.min_frequency_of_ram AND ram.clock_frequency <= cpu.max_frequency_of_ram;
INSERT INTO motherboard_ram(motherboard_id, ram_id) SELECT motherboard.motherboardid, ram.ramid FROM motherboard, ram WHERE position(ram.ddr in motherboard.supported_ddr) > 0 AND ram.clock_frequency >= motherboard.min_frequency_of_ram AND ram.clock_frequency <= motherboard.max_frequency_of_ram AND ram.form_factor = motherboard.form_factor_of_ram;
INSERT INTO cpu_motherboard(cpu_id, motherboard_id) SELECT cpu.cpuid, motherboard.motherboardid FROM cpu, motherboard WHERE (string_to_array(motherboard.supported_ddr, ', ') && string_to_array(cpu.supported_ddr, ', ')) AND motherboard.max_frequency_of_ram >= cpu.min_frequency_of_ram AND motherboard.min_frequency_of_ram <= cpu.max_frequency_of_ram;
--------------------------------------------------данные для нагрузочного тестирования
-- ALTER TABLE cpu_data ADD COLUMN temp_id SERIAL;
-- ALTER TABLE ram_data ADD COLUMN temp_id SERIAL;
-- ALTER TABLE cooling_data ADD COLUMN temp_id SERIAL;
-- ALTER TABLE motherboard_data ADD COLUMN temp_id SERIAL;
-- ALTER TABLE powerbox_data ADD COLUMN temp_id SERIAL;
-- ALTER TABLE videocard_data ADD COLUMN temp_id SERIAL;


-- INSERT INTO cpu (name, brand, series, model, supported_ddr, price)
-- SELECT c2.name, c2.brand, c2.series, c2.model, c2.supported_ddr, c2.price
-- FROM cpu_data c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM cpu c1
--     WHERE c1.cpuid = c2.temp_id
--     AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.series = c2.series
--       AND c1.model = c2.model
--       AND c1.supported_ddr = c2.supported_ddr
--       AND c1.price = c2.price
-- );
-- INSERT INTO cpu (name, brand, series, model, basic_frequency, supported_ddr, min_frequency_of_ram, max_frequency_of_ram, price)
-- SELECT c2.name, c2.brand, c2.series, c2.model, c2.basic_frequency, c2.supported_ddr, c2.min_frequency_of_ram, c2.max_frequency_of_ram, c2.price
-- FROM cpu_data_1 c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM cpu c1
--     WHERE c1.cpuid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.series = c2.series
--       AND c1.model = c2.model
--       AND c1.basic_frequency = c2.basic_frequency
--       AND c1.supported_ddr = c2.supported_ddr
--       AND c1.min_frequency_of_ram = c2.min_frequency_of_ram
--       AND c1.max_frequency_of_ram = c2.max_frequency_of_ram
--       AND c1.price = c2.price
-- );


--
-- INSERT INTO ram (name, brand, rammemory, ddr, price)
-- SELECT c2.name, c2.brand, c2.rammemory, c2.ddr, c2.price
-- FROM ram_data c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM ram c1
--     WHERE c1.ramid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.rammemory = c2.rammemory
--       AND c1.ddr = c2.ddr
--       AND c1.price = c2.price
-- );

-- INSERT INTO ram (name, brand, rammemory, ddr, form_factor, clock_frequency, price)
-- SELECT c2.name, c2.brand, c2.rammemory, c2.ddr, c2.form_factor, c2.clock_frequency, c2.price
-- FROM ram_data_1 c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM ram c1
--     WHERE c1.ramid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.rammemory = c2.rammemory
--       AND c1.ddr = c2.ddr
--       AND c1.form_factor = c2.form_factor
--       AND c1.clock_frequency = c2.clock_frequency
--       AND c1.price = c2.price
-- );

-- INSERT INTO cooling (name, brand, max_speed, count_ventilators, price)
-- SELECT c2.name, c2.brand, c2.max_speed, c2.count_ventilators, c2.price
-- FROM cooling_data c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM cooling c1
--     WHERE c1.coolingid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.max_speed = c2.max_speed
--       AND c1.count_ventilators = c2.count_ventilators
--       AND c1.price = c2.price
-- );

-- INSERT INTO motherboard (name, brand, count_slots, price)
-- SELECT c2.name, c2.brand, c2.count_slots, c2.price
-- FROM motherboard_data c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM motherboard c1
--     WHERE c1.motherboardid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.count_slots = c2.count_slots
--       AND c1.price = c2.price
-- );

-- INSERT INTO motherboard (name, brand, count_slots, supported_ddr, form_factor_of_ram, min_frequency_of_ram, max_frequency_of_ram, price)
-- SELECT c2.name, c2.brand, c2.count_slots, c2. supported_ddr, c2.form_factor_of_ram, c2.min_frequency_of_ram, c2.max_frequency_of_ram, c2.price
-- FROM motherboard_data_1 c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM motherboard c1
--     WHERE c1.motherboardid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.count_slots = c2.count_slots
--       AND c1.supported_ddr = c2.supported_ddr
--       AND c1.form_factor_of_ram = c2.form_factor_of_ram
--       AND c1.min_frequency_of_ram = c2.min_frequency_of_ram
--       AND c1.max_frequency_of_ram = c2.max_frequency_of_ram
--       AND c1.price = c2.price
-- );

-- INSERT INTO powerbox (name, brand, power, form_factor, price)
-- SELECT c2.name, c2.brand, c2.power, c2.form_factor, c2.price
-- FROM powerbox_data c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM powerbox c1
--     WHERE c1.powerboxid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.power = c2.power
--       AND c1.form_factor = c2.form_factor
--       AND c1.price = c2.price
-- );

-- INSERT INTO videocard (name, brand, series, vmemory, price)
-- SELECT c2.name, c2.brand, c2.series, c2.vmemory, c2.price
-- FROM videocard_data c2
-- WHERE NOT EXISTS (
--     SELECT 1
--     FROM videocard c1
--     WHERE c1.videocardid = c2.temp_id
--       AND c1.name = c2.name
--       AND c1.brand = c2.brand
--       AND c1.series = c2.series
--       AND c1.vmemory = c2.vmemory
--       AND c1.price = c2.price
-- );
--------------------------------------------------
-- drop view if exists cpu_view;
-- drop view if exists cooling_view;
-- drop view if exists motherboard_view;
-- drop view if exists powerbox_view;
-- drop view if exists ram_view;
-- drop view if exists videocard_view;