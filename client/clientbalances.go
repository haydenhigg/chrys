package client

func (c *Client) GetBalances() (map[string]float64, error) {
	// check store
	if balances, ok := c.Store.TryGetBalances(); ok {
		return balances, nil
	}

	// retrieve from data source
	balances, err := c.Connector.GetBalances()
	if err != nil {
		return nil, err
	}

	c.Store.Balances = balances

	return balances, nil
}
