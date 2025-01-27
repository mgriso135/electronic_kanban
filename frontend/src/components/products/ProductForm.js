import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const ProductForm = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [productID, setProductID] = useState('');
    const [name, setName] = useState('');
    const [isEdit, setIsEdit] = useState(false);


    useEffect(() => {
        if (id) {
            setIsEdit(true);
            const fetchProduct = async () => {
              try {
                const response = await api.get(`/products/${id}`);
                setProductID(response.data.product_id);
                setName(response.data.name);
              } catch (error) {
                console.error("Error fetching product:", error)
              }
            };
            fetchProduct();
        }
    }, [id]);


    const handleSubmit = async (e) => {
      e.preventDefault();
      const productData = {
        product_id: productID,
        name: name
      };
      try {
          if(isEdit){
            await api.put(`/products/${id}`, productData);
          } else {
              await api.post('/products', productData);
          }
          navigate('/products');
      } catch (error) {
        console.error('Error saving product:', error);
      }
    };


    return (
        <div>
            <h2>{isEdit ? 'Edit Product' : 'Create New Product'}</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Product ID:</label>
                    <input type="text" value={productID} onChange={(e) => setProductID(e.target.value)} disabled={isEdit} />
                </div>
                <div>
                    <label>Name:</label>
                    <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
                </div>
                <button type="submit">{isEdit ? 'Update Product' : 'Create Product'}</button>
                <Link to="/products">Cancel</Link>
            </form>
        </div>
    );
};

export default ProductForm;