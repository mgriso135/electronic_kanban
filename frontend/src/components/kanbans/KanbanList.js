import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

const KanbanList = () => {
    const [kanbans, setKanbans] = useState([]);
    const [products, setProducts] = useState([]);
    const [selectedProductIds, setSelectedProductIds] = useState([]);

    useEffect(() => {
        fetchProducts();
    }, []);

    const fetchKanbans = useCallback(async () => {
        try {
            const params = new URLSearchParams();
            selectedProductIds.forEach(id => params.append("product_id", id));

            const response = await api.get('/kanbans?' + params.toString());
            setKanbans(response.data);
        } catch (error) {
            console.error('Error fetching kanbans:', error);
        }
    }, [selectedProductIds]);

    const fetchProducts = async () => {
        try {
            const response = await api.get('/products');
            setProducts(response.data);
        } catch (error) {
            console.error('Error fetching products:', error);
        }
    };

    const handleProductFilterChange = (event) => {
        const selectedOptions = Array.from(event.target.selectedOptions, option => option.value);
        setSelectedProductIds(selectedOptions);
    };

    useEffect(() => {
        fetchKanbans();
    }, [fetchKanbans]);


    const handleDelete = async (id) => {
        if (window.confirm("Are you sure you want to delete this kanban?")) {
            try {
                await api.delete(`/kanbans/${id}`);
                setKanbans(kanbans.filter(kanban => kanban.id !== id));
            } catch (error) {
                console.error('Error deleting kanban:', error);
            }
        }
    };


    return (
        <div>
            <h2>Kanban List</h2>
            <Link to="/kanbans/new">Create New Kanban</Link> {/* Add "Create New Kanban" Link */}


            <div>
                <label htmlFor="productFilter">Filter by Product:</label>
                <select
                    id="productFilter"
                    multiple
                    value={selectedProductIds}
                    onChange={handleProductFilterChange}
                >
                    {products.map(product => (
                        <option key={product.product_id} value={product.product_id}>{product.name} ({product.product_id})</option>
                    ))}
                </select>
            </div>


            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Data Aggiornamento</th>
                        <th>Lead Time</th>
                        <th>Container Type</th>
                        <th>Quantity</th>
                        <th>Status Current</th>
                        <th>Product ID</th>
                        <th>Product Name</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {kanbans.map(kanban => (
                        <tr key={kanban.id}>
                            <td>{kanban.id}</td>
                            <td>{new Date(kanban.data_aggiornamento).toLocaleString()}</td>
                            <td>{kanban.leadtime_days}</td>
                            <td>{kanban.tipo_contenitore}</td>
                            <td>{kanban.quantity}</td>
                            <td>
                                <span style={{ 
                                    backgroundColor: kanban.status_color || '#e0e0e0', 
                                    color: '#fff', 
                                    padding: '4px 8px', 
                                    borderRadius: '12px', 
                                    fontSize: '0.85em',
                                    fontWeight: 'bold',
                                    textShadow: '0 1px 2px rgba(0,0,0,0.2)'
                                }}>
                                    {kanban.status_name || `ID: ${kanban.status_current}`}
                                </span>
                            </td>
                            <td>{kanban.product_id}</td>
                            <td>{kanban.product_name}</td>
                            <td>
                                <Link to={`/kanbans/${kanban.id}/edit`}>Edit</Link> {/* Add "Edit" Link */}
                                <button onClick={() => handleDelete(kanban.id)}>Delete</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default KanbanList;