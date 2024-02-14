INSERT into coffee (coffee_name, coffee_description, coffee_price) VALUES 
    ('coffee', 'the original - simple, bold and satisfying', 2.25),
    ('espresso', 'a concentrated and small super dose of coffee', 2.25),
    ('latte', 'espresso with a little bit of steamed milk', 3.00),
    ('cappuchino', 'espresso and a small bit of extra frothy milk', 3.00),
    ('americano', 'espresso and extra water', 2.50)
    ON CONFLICT DO NOTHING;