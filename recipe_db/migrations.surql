DEFINE TABLE ingredients SCHEMAFULL PERMISSIONS NONE;
DEFINE TABLE recipe SCHEMALESS PERMISSIONS NONE;
DEFINE TABLE recipe_step SCHEMAFULL PERMISSIONS NONE;
DEFINE TABLE recipe_tag SCHEMALESS PERMISSIONS NONE;

DEFINE FIELD name ON recipe TYPE string ASSERT $value != NONE;
DEFINE FIELD ingredients ON recipe TYPE array;
DEFINE FIELD instructions ON recipe TYPE array;
DEFINE FIELD picture_url ON recipe TYPE string ASSERT $value != NONE;

DEFINE FIELD name ON ingredients TYPE string ASSERT $value != NONE;
DEFINE FIELD description ON ingredients TYPE string ASSERT $value != NONE;

DEFINE FIELD name ON recipe_tag TYPE string ASSERT $value != NONE;

CREATE recipe_tag:fgcrl555snthd CONTENT {
    name: "Dinner"
};

CREATE recipe_tag:fgcrl555dhtns CONTENT {
    name: "Lunch"
};

CREATE recipe_tag:fgcrl555vwmbx CONTENT {
    name: "Breakfast"
};

CREATE recipe_tag:lrcgf555snthd CONTENT {
    name: "Savory"
};

CREATE recipe_tag:lrcgf555snthd CONTENT {
    name: "Spicy"
};

CREATE recipe_tag:aoeui626snthd CONTENT {
    name: "Sweet"
};

CREATE recipe_tag:lrc44555sntao CONTENT {
    name: "Umami"
};

CREATE recipe_tag:00cge9564nt0d CONTENT {
    name: "Desert"
};

CREATE ingredients:00b2pg847d7b8r08t08t CONTENT {
    name: "Eggs",
    description: "Large eggs"
};

CREATE ingredients:123abc456xyz CONTENT {
    name: "Flour",
    description: "White all-purpose flour"
};

CREATE ingredients:456def789uvw CONTENT {
    name: "Sugar",
    description: "Granulated sugar"
};

CREATE ingredients:789ghi123lmn CONTENT {
    name: "Milk",
    description: "Whole milk"
};

CREATE ingredients:abc456def789 CONTENT {
    name: "Butter",
    description: "Unsalted butter"
};

// Create ingredients entries
CREATE ingredients:uvw123xyz456 CONTENT {
    name: "Salt",
    description: "Table salt"
};

CREATE ingredients:lmn789pqr012 CONTENT {
    name: "Baking Powder",
    description: "Double-acting baking powder"
};

CREATE ingredients:pqr456stu789 CONTENT {
    name: "Vanilla Extract",
    description: "Pure vanilla extract"
};

CREATE ingredients:xyz789abc012 CONTENT {
    name: "Chocolate Chips",
    description: "Semi-sweet chocolate chips"
};

CREATE ingredients:def012ghi345 CONTENT {
    name: "Brown Sugar",
    description: "Light brown sugar"
};

CREATE ingredients:stu456vwx789 CONTENT {
    name: "Cinnamon",
    description: "Ground cinnamon"
};

CREATE ingredients:ghi123jkl234 CONTENT {
    name: "Olive Oil",
    description: "Extra virgin olive oil"
};

CREATE ingredients:jkl567mno890 CONTENT {
    name: "Onions",
    description: "Yellow onions"
};

CREATE ingredients:vwx789yza012 CONTENT {
    name: "Garlic",
    description: "Fresh garlic"
};

CREATE ingredients:klm234nop567 CONTENT {
    name: "Tomatoes",
    description: "Roma tomatoes"
};

CREATE ingredients:opq567rst890 CONTENT {
    name: "Basil",
    description: "Fresh basil"
};

CREATE ingredients:aoeu555snthd CONTENT {
    name: "Bacon",
    description: "Preferably thick cut, or any thick cut fatty meat."
};

CREATE ingredients:snth444aoeui CONTENT {
    name: "Tomatoes",
    description: "Peferably San Marzano Tomatoes Peeled."
};

CREATE ingredients:aoeui555snth CONTENT {
    name: "Spaghetti",
    description: "Any type will do",
};

CREATE ingredients:snthd555aoeu CONTENT {
    name: "Oregano",
    description: "Any store brand dried Oregano will do; however, fresh oregano is the best!"
};

-- Create recipes
CREATE recipe:abc123xyz456 CONTENT {
    name: "Classic Pancakes",
    ingredients: [
        ingredients:00b2pg847d7b8r08t08t,
        ingredients:123abc456xyz,
        ingredients:456def789uvw,
        ingredients:789ghi123lmn,
        ingredients:abc456def789,
        ingredients:uvw123xyz456,
        ingredients:lmn789pqr012,
        ingredients:pqr456stu789,
        ingredients:xyz789abc012
    ],
    instructions: [
        {
            step: "Mix dry ingredients",
            type: "prep",
            image_file_loc: "images/step1.jpg"
        },
        {
            step: "Whisk wet ingredients",
            type: "prep",
            image_file_loc: "images/step2.jpg"
        },
        {
            step: "Combine wet and dry ingredients",
            type: "prep",
            image_file_loc: "images/step3.jpg"
        },
        {
            step: "Cook pancakes on griddle",
            type: "cooking",
            image_file_loc: "images/step4.jpg"
        },
        {
            step: "Flip pancakes when bubbles form",
            type: "cooking",
            image_file_loc: "images/step5.jpg"
        },
        {
            step: "Serve with syrup",
            type: "plating",
            image_file_loc: "images/step6.jpg"
        }
    ],
    picture_url: "images/pancakes.jpg",
    tags: [
        recipe_tag:fgcrl555snthd,
        recipe_tag:aoeui626snthd
    ]
};

CREATE recipe:def789ghi123 CONTENT {
    name: "Spaghetti Carbonara",
    ingredients: [
        ingredients:00b2pg847d7b8r08t08t,
        ingredients:aoeui555snth,
        ingredients:aoeui626snthd,
        ingredients:snth444aoeui,
        ingredients:aoeui555snth,
        ingredients:aoeui626snthd,
        ingredients:snth444aoeui
    ],
    instructions: [
        {
            step: "Boil spaghetti",
            type: "prep",
            image_file_loc: "images/step1.jpg"
        },
        {
            step: "Cook bacon until crispy",
            type: "cooking",
            image_file_loc: "images/step2.jpg"
        },
        {
            step: "Whisk eggs, cheese, and pepper",
            type: "prep",
            image_file_loc: "images/step3.jpg"
        },
        {
            step: "Mix cooked spaghetti and bacon",
            type: "cooking",
            image_file_loc: "images/step4.jpg"
        },
        {
            step: "Stir in egg mixture",
            type: "cooking",
            image_file_loc: "images/step5.jpg"
        },
        {
            step: "Serve with extra cheese and pepper",
            type: "plating",
            image_file_loc: "images/step6.jpg"
        }
    ],
    picture_url: "images/carbonara.jpg",
    tags: [
        recipe_tag:fgcrl555snthd,
        recipe_tag:lrcgf555snthd
    ]
};
