INSERT into coffee (coffee_name, coffee_description, coffee_price, coffee_caffeine, coffee_calories) VALUES 
    ('coffee', 'the original - simple, bold and satisfying', 2.25, '90mg',5),
    ('espresso', 'a concentrated and small super dose of coffee', 2.25,'120mg',5),
    ('latte', 'espresso with a little bit of steamed milk', 4.50, '120mg',100),
    ('cappuchino', 'espresso and a small bit of extra frothy milk', 4.50, '120mg', 100),
    ('americano', 'espresso and extra water', 2.50, '120mg',5)
    ON CONFLICT DO NOTHING;