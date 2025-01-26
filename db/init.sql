
-- Table: attachment_figures
CREATE TABLE attachment_figures (
                                    id SERIAL PRIMARY KEY,
                                    name VARCHAR(255),
                                    description TEXT,
                                    image TEXT, -- File location
                                    birth_date DATE,
                                    death_date DATE,
                                    gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')) -- 'O' for Other
);

-- Table: family_descendant
CREATE TABLE family_descendant (
                                   id SERIAL PRIMARY KEY,
                                   name VARCHAR(255),
                                   description TEXT,
                                   image TEXT, -- File location
                                   family_parent INT REFERENCES family_descendant(id) ON DELETE SET NULL,
                                   attachment_parent INT REFERENCES attachment_figures(id) ON DELETE SET NULL,
                                   generation INT,
                                   gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')), -- 'O' for Other
                                   birth_date DATE,
                                   death_date DATE
);



-- Join Table: descendant_attachments
CREATE TABLE descendant_attachments (
                                        id SERIAL PRIMARY KEY,
                                        attachment_figure_id INT NOT NULL REFERENCES attachment_figures(id) ON DELETE CASCADE,
                                        descendant_id INT NOT NULL REFERENCES family_descendant(id) ON DELETE CASCADE
);

-- Indexes for faster lookups
CREATE INDEX idx_family_descendant_family_parent ON family_descendant(family_parent);
CREATE INDEX idx_family_descendant_attachment_parent ON family_descendant(attachment_parent);
CREATE INDEX idx_descendant_attachments_attachment ON descendant_attachments(attachment_figure_id);
CREATE INDEX idx_descendant_attachments_descendant ON descendant_attachments(descendant_id);





-- insertions

-- Insert Adam (First Generation)
INSERT INTO family_descendant (name, description, image, generation, gender, birth_date, death_date)
VALUES ('Adam', 'The first man in the family tree.', '/images/adam.png', 1, 'M', '1900-01-01', '1980-01-01');

-- Insert Eve (Attachment Figure)
INSERT INTO attachment_figures (name, description, image, birth_date, death_date, gender)
VALUES ('Eve', 'The first woman, related to Adam by marriage.', '/images/eve.png', '1905-01-01', '1990-01-01', 'F');

-- Insert children of Adam and Eve
INSERT INTO family_descendant (name, description, image, family_parent, attachment_parent, generation, gender, birth_date, death_date)
VALUES
    ('Cain', 'First child of Adam and Eve.', '/images/cain.png', 1, 1, 2, 'M', '1925-01-01', '2000-01-01'),
    ('Abel', 'Second child of Adam and Eve.', '/images/abel.png', 1, 1, 2, 'M', '1928-01-01', '1950-01-01'),
    ('Seth', 'Third child of Adam and Eve.', '/images/seth.png', 1, 1, 2, 'M', '1930-01-01', '2005-01-01');

-- Insert grandchildren (children of Cain and his partner)
INSERT INTO family_descendant (name, description, image, family_parent, generation, gender, birth_date, death_date)
VALUES
    ('Enoch', 'Child of Cain.', '/images/enoch.png', 2, 3, 'M', '1950-01-01', '2020-01-01'),
    ('Irad', 'Another child of Cain.', '/images/irad.png', 2, 3, 'M', '1955-01-01', '2025-01-01');

-- Insert a descendant attachment figure not directly related by marriage (example: a family friend)
INSERT INTO attachment_figures (name, description, image, birth_date, death_date, gender)
VALUES ('Family Friend', 'A close family friend who helped raise the children.', '/images/family_friend.png', '1915-01-01', '1995-01-01', 'O');

-- Link Adam's children to the family friend in the descendant_attachments table
INSERT INTO descendant_attachments (attachment_figure_id, descendant_id)
VALUES
    (2, 2), -- Family Friend attached to Cain
    (2, 3), -- Family Friend attached to Abel
    (2, 4); -- Family Friend attached to Seth


/*

select * from family_descendant;
-- get family parents name
SELECT
    fd.name AS descendant_name,
    fp.name AS family_parent_name,
    af.name AS attachment_parent_name
FROM
    family_descendant fd
        LEFT JOIN
    family_descendant fp ON fd.family_parent = fp.id
        LEFT JOIN
    attachment_figures af ON fd.attachment_parent = af.id;



-- recursively get all parents
WITH RECURSIVE parent_tree AS (
    -- Base case: Start with the descendant at generation 1
    SELECT
        fd.id AS descendant_id,
        fd.name AS descendant_name,
        fd.family_parent AS family_parent_id,
        fd.attachment_parent AS attachment_parent_id,
        1 AS generation -- Start with generation 1 for the descendant
    FROM
        family_descendant fd
    WHERE
        fd.id = 6  -- Replace this with the descendant's id

    UNION ALL

    -- Recursive case: Get parents and increment generation based on the parent's generation
    SELECT
        fd.id AS descendant_id,
        fd.name AS descendant_name,
        fd.family_parent AS family_parent_id,
        fd.attachment_parent AS attachment_parent_id,
        pt.generation + 1 AS generation -- Increment generation upwards by 1 for each parent
    FROM
        family_descendant fd
            JOIN
        parent_tree pt ON fd.id = pt.family_parent_id OR fd.id = pt.attachment_parent_id
    WHERE
        pt.generation < 10  -- Prevent infinite recursion if there are loops (optional)
)
SELECT
    pt.descendant_name,
    MAX(pt.generation) OVER () - pt.generation +1 AS generation,
        fp.name AS family_parent_name,
    af.name AS attachment_parent_name
FROM
    parent_tree pt
        LEFT JOIN
    family_descendant fp ON pt.family_parent_id = fp.id
        LEFT JOIN
    attachment_figures af ON pt.attachment_parent_id = af.id
ORDER BY
    generation ;  -- Order by generation in ascending order (starting from 1)
*/