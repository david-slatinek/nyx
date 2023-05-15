USE recommend_db;

SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS recommend;
DROP TABLE IF EXISTS recommend_follow;
DROP TABLE IF EXISTS order_table;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE recommend
(
    id_recommend     INT          NOT NULL AUTO_INCREMENT,
    user_id          VARCHAR(255) NOT NULL,
    fk_main_category VARCHAR(255) NOT NULL,
    fk_sub_category  VARCHAR(255) NOT NULL,
    category_name    VARCHAR(255) NOT NULL,
    score            FLOAT        NOT NULL,
    fk_main_dialog   VARCHAR(255) NOT NULL,
    fk_dialog        VARCHAR(255) NOT NULL,
    recommended_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_recommend)
);

CREATE TABLE recommend_follow
(
    id_recommend_follow INT      NOT NULL AUTO_INCREMENT,
    fk_recommend        INT      NOT NULL,
    click_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_recommend_follow),
    FOREIGN KEY (fk_recommend) REFERENCES recommend (id_recommend)
);

CREATE TABLE order_table
(
    id_order_table INT      NOT NULL AUTO_INCREMENT,
    fk_recommend   INT      NOT NULL,
    order_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    quantity       INT      NOT NULL,
    PRIMARY KEY (id_order_table),
    FOREIGN KEY (fk_recommend) REFERENCES recommend (id_recommend)
);
