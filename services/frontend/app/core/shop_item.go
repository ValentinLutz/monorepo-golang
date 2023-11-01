package core

type ShopItem struct {
	SKU         string
	Name        string
	Description string
	Price       float64
	ImageName   string
	ImageAlt    string
}

var ShopItems = []ShopItem{
	{
		SKU:         "FRUIT-0001",
		Name:        "Apple",
		Description: "Discover the essence of freshness with our premium apples. Crisp, juicy, and packed with flavor, these apples are nature's perfection in every bite. Whether for snacking or enhancing your recipes, our apples are the epitome of quality and taste.",
		Price:       2.49,
		ImageName:   "apple.png",
		ImageAlt:    "apple",
	},
	{
		SKU:         "FRUIT-0002",
		Name:        "Banana",
		Description: "Bananas are one of the most popular fruits worldwide. They contain essential nutrients that can have a protective impact on health. Eating bananas can help lower blood pressure and may reduce the risk of cancer.",
		Price:       1.57,
		ImageName:   "banana.png",
		ImageAlt:    "banana",
	},
	{
		SKU:         "FRUIT-0003",
		Name:        "Orange",
		Description: "Oranges are among the world's most popular fruits, as they're both tasty and nutritious. They are a good source of vitamin C, as well as several other vitamins, minerals, and antioxidants. For this reason, they may lower your risk of heart disease and kidney stones.",
		Price:       2.19,
		ImageName:   "orange.png",
		ImageAlt:    "orange",
	},
	{
		SKU:         "FRUIT-0004",
		Name:        "Kiwi",
		Description: "Kiwis are small fruits that pack a lot of flavor and plenty of health benefits. Their green flesh is sweet and tangy. It's also full of nutrients like vitamin C, vitamin K, vitamin E, folate, and potassium. They also have a lot of antioxidants and are a good source of fiber.",
		Price:       1.58,
		ImageName:   "kiwi.png",
		ImageAlt:    "kiwi",
	},
	{
		SKU:         "FRUIT-0005",
		Name:        "Pineapple",
		Description: "Pineapple is a delicious tropical fruit, celebrated for centuries, not only for its unique taste but also for its miraculous health benefits. The health and medicinal benefits of pineapple include boosting the immune system, and respiratory health, aiding digestion, and strengthening bones.",
		Price:       3.99,
		ImageName:   "pineapple.png",
		ImageAlt:    "pineapple",
	},
	{
		SKU:         "FRUIT-0006",
		Name:        "Strawberry",
		Description: "Strawberries are bright red, juicy, and sweet. They're an excellent source of vitamin C and manganese and also contain decent amounts of folate (vitamin B9) and potassium. Strawberries are very rich in antioxidants and plant compounds, which may have benefits for heart health and blood sugar control.",
		Price:       2.99,
		ImageName:   "strawberry.png",
		ImageAlt:    "strawberry",
	},
	{
		SKU:         "FRUIT-0007",
		Name:        "Lime",
		Description: "Limes are a good source of magnesium and potassium, which promote heart health. Potassium can naturally lower blood pressure and improve blood circulation, which reduces your risk of a heart attack and stroke. Research shows that limes may promote healthy cell function and reduce your risk of cancer.",
		Price:       1.99,
		ImageName:   "lime.png",
		ImageAlt:    "lime",
	},
	{
		SKU:         "FRUIT-0008",
		Name:        "Cherry",
		Description: "Cherries are a good source of fiber, vitamins, and minerals, including potassium, calcium, vitamin A, and folic acid. They are also well known for their antioxidant properties. A growing body of evidence suggests that cherries may help reduce inflammation and lower the risk of gout, diabetes, heart disease, and certain types of cancer.",
		Price:       2.99,
		ImageName:   "cherry.png",
		ImageAlt:    "cherry",
	},
	{
		SKU:         "FRUIT-0009",
		Name:        "Raspberry",
		Description: "Raspberries are low in calories but high in fiber, vitamins, minerals, and antioxidants. They may protect against diabetes, cancer, obesity, arthritis, and other conditions and may even provide anti-aging effects. Raspberries are easy to add to your diet and make a tasty addition to breakfast, lunch, dinner, or dessert.",
		Price:       3.99,
		ImageName:   "raspberry.png",
		ImageAlt:    "raspberry",
	},
}
