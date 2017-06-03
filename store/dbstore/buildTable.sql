CREATE TABLE `ballon_status` (
    `team_key` TEXT NOT NULL,
    `problem_index` INTEGER NOT NULL,
    `solution_time` INTEGER NULL,
    `team_name` TEXT NOT NULL,
    `seat_id` TEXT NULL,
    `status` TEXT NULL,
    `is_marked` INTEGER NULL,
    PRIMARY KEY (`team_key`, `problem_index`)
);
CREATE TABLE `kv` (
    `key` TEXT PRIMARY KEY NOT NULL,
    `value` BLOB NULL
);
CREATE TABLE `user` (
    `account` TEXT PRIMARY KEY NOT NULL,
    `password` TEXT NOT NULL,
    `role` TEXT NOT NULL DEFAULT 'normal',
    `display_name` TEXT NULL,
    `nick_name` TEXT NULL,
    `school` TEXT NULL,
    `is_star` TEXT NULL,
    `is_girl` TEXT NULL,
    `seat_id` TEXT NULL,
    `coach` TEXT NULL,
    `player1` TEXT NULL,
    `player2` TEXT NULL,
    `player3` TEXT NULL,
    `site` TEXT,
    `team_key` TEXT UNIQUE
);

INSERT INTO `user` (`account`, `password`, `role`) VALUES('admin', 'ElPsyCongroo', 'admin');


-- CREATE TABLE `problem` (
--     `id` INTEGER PRIMARY KEY NOT NULL,
--     `number_solved` INTEGER NULL,
--     `best_solution_time` INTEGER NULL,
--     `last_solution_time` INTEGER NULL,
--     `title` TEXT NULL,
--     `short_name` TEXT NULL
-- );
-- CREATE TABLE `team` (
--     `team_key` TEXT PRIMARY KEY NOT NULL,
--     `team_name` TEXT NOT NULL,
--     `points` INTEGER NOT NULL,
--     `problems_attempted` INTEGER NOT NULL,
--     `rank` INTEGER NOT NULL,
--     `solved` INTEGER NOT NULL,
--     `total_attempts` INTEGER NOT NULL
-- );
-- CREATE TABLE `standings_header` (
--     `current_date` TEXT NOT NULL,
--     `problem_count` INTEGER NOT NULL,
--     `problems_attempted` INTEGER NOT NULL,
--     `total_attempts` INTEGER NOT NULL,
--     `total_solved` INTEGER NOT NULL
-- );
-- CREATE TABLE `team_problem` (
--     `team_key` TEXT NOT NULL,
--     `problem_id` INTEGER NOT NULL,
--     `attempts` INTEGER NULL,
--     `is_pending` INTEGER NULL,
--     `is_solved` INTEGER NULL,
--     `points` INTEGER NULL,
--     `solution_time` INTEGER NULL,
--     `is_first_solved` INTEGER NULL
-- );
-- CREATE TABLE `users` (
--     `account` TEXT PRIMARY KEY NOT NULL,
--     `password` TEXT NOT NULL,
--     `role` TEXT NOT NULL DEFAULT 'normal',
--     `display_name` TEXT NULL,
--     `nick_name` TEXT NULL,
--     `school` TEXT NULL,
--     `is_star` TEXT NULL,
--     `is_girl` TEXT NULL,
--     `seat_id` TEXT NULL,
--     `coach` TEXT NULL,
--     `player1` TEXT NULL,
--     `player2` TEXT NULL,
--     `player3` TEXT NULL
-- );

-- INSERT INTO `users` (`account`, `password`, `role`) VALUES('admin', 'ElPsyCongroo', 'admin');

-- CREATE TABLE `contest_standing` (
--     `x_m_l_name` TEXT NULL,
--     `standings_header` TEXT NULL,
--     `team_standings` TEXT NULL
-- );
-- CREATE TABLE `problem` (
--     `i_d` INTEGER PRIMARY KEY NOT NULL,
--     `number_solved` INTEGER NULL,
--     `best_solution_time` INTEGER NULL,
--     `last_solution_time` INTEGER NULL,
--     `title` TEXT NULL,
--     `short_name` TEXT NULL
-- );
-- CREATE TABLE `standings_header` (
--     `current_date` TEXT NULL,
--     `problem_count` INTEGER NULL,
--     `problems_attempted` INTEGER NULL,
--     `total_attempts` INTEGER NULL,
--     `total_solved` INTEGER NULL,
--     `problems` TEXT NULL
-- );
-- CREATE TABLE `problem_summary_info` (
--     `index` INTEGER NULL,
--     `attempts` INTEGER NULL,
--     `is_pending` INTEGER NULL,
--     `is_solved` INTEGER NULL,
--     `points` INTEGER NULL,
--     `solution_time` INTEGER NULL,
--     `is_first_solved` INTEGER NULL
-- );
-- CREATE TABLE `team_standing` (
--     `first_solved` INTEGER NULL,
--     `index` INTEGER NULL,
--     `last_solved` INTEGER NULL,
--     `points` INTEGER NULL,
--     `problems_attempted` INTEGER NULL,
--     `rank` INTEGER NULL,
--     `solved` INTEGER NULL,
--     `team_name` TEXT NULL,
--     `total_attempts` INTEGER NULL,
--     `problem_summary_infos` TEXT NULL
-- );
