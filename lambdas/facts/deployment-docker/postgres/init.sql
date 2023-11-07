CREATE TABLE IF NOT EXISTS public.fact
(
    fact_id   INT GENERATED ALWAYS AS IDENTITY,
    fact_text VARCHAR,
    PRIMARY KEY (fact_id)
);

INSERT INTO public.fact
    (fact_text)
VALUES ('The Eiffel Tower has 1,665 steps.'),
       ('The average person will spend six months of their life waiting for red lights to turn green.'),
       ('Honey never spoils. Archaeologists have found pots of honey in ancient Egyptian tombs that are over 3,000 years old and still perfectly edible.'),
       ('The world''s largest desert is not the Sahara but Antarctica.'),
       ('The human brain is the most energy-consuming organ, using up to 20% of the body''s total energy.');
