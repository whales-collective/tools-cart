package tools

import "github.com/openai/openai-go"

func GetToolsCatalog() []openai.ChatCompletionToolParam {

	searchProducts := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "search_products",
			Description: openai.String("Search for products by query, category, or price range"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query for product name or description",
					},
					"category": map[string]interface{}{
						"type":        "string",
						"description": "Product category (electronics, clothing, books, home, sports, beauty, toys, food)",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results to return (default: 10)",
					},
				},
			},
		},
	}

	addToCart := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "add_to_cart",
			Description: openai.String("Add a quantity of a product to the shopping cart"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"product_name": map[string]interface{}{
						"type":        "string",
						"description": "The name of the product to add",
					},
					"quantity": map[string]interface{}{
						"type":        "integer",
						"description": "Quantity to add (default: 1)",
					},
				},
				"required": []string{"product_name"},
			},
		},
	}

	removeFromCart := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "remove_from_cart",
			Description: openai.String("Remove a product from the shopping cart"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"product_name": map[string]interface{}{
						"type":        "string",
						"description": "The name of the product to remove",
					},
				},
				"required": []string{"product_name"},
			},
		},
	}

	viewCart := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "view_cart",
			Description: openai.String("View the current shopping cart contents and totals"),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	updateQuantity := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "update_quantity",
			Description: openai.String("Update the quantity of a product in the cart"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"product_name": map[string]interface{}{
						"type":        "string",
						"description": "The name of the product to update",
					},
					"quantity": map[string]interface{}{
						"type":        "integer",
						"description": "New quantity (use 0 to remove)",
					},
				},
				"required": []string{"product_name", "quantity"},
			},
		},
	}

	checkOut := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "checkout",
			Description: openai.String("Process checkout for the current cart"),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	tools := []openai.ChatCompletionToolParam{
		searchProducts,
		addToCart,
		removeFromCart,
		viewCart,
		updateQuantity,
		checkOut,
	}
	return tools
}
