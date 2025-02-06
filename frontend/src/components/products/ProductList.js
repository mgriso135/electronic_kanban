import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

const ProductList = () => {
  const [products, setProducts] = useState(null);

    useEffect(() => {
        const fetchProducts = async () => {
            try {
              const response = await api.get('/products');
              console.log("API response:", response);
                setProducts(response.data);
            } catch (error) {
                console.error('Error fetching products:', error);
              setProducts([]);
            }
        };
        fetchProducts();
    }, []);

  const handleDelete = async (id) => {
      try {
        await api.delete(`/products/${id}`);
        setProducts(products.filter(product => product.product_id !== id));
      } catch (error) {
        console.error('Error deleting product:', error);
      }
  }

    return (
      <div>
          <h2>Products</h2>
          <Link to="/products/new">Create New Product</Link>
          {products && products.length > 0 ? (
              <table>
                  <thead>
                    <tr>
                          <th>Product ID</th>
                          <th>Name</th>
                          <th>Actions</th>
                    </tr>
                  </thead>
                <tbody>
                  {products.map(product => (
                      <tr key={product.product_id}>
                          <td>{product.product_id}</td>
                          <td>{product.name}</td>
                        <td>
                           <Link to={`/products/${product.product_id}/edit`}>Edit</Link>
                           <button onClick={() => handleDelete(product.product_id)}>Delete</button>
                         </td>
                      </tr>
                  ))}
                </tbody>
              </table>
            ) : (
            <p>No products found</p>
            )}
      </div>
    );
};

export default ProductList;