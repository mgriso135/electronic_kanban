import React from 'react';
import api from '../../services/api';

const KanbanCard = ({ kanban, dashboardType, setKanbans, productID }) => { // Added productID prop

    const handleStatusChange = async () => {
        try {
            const response = await api.put(`/kanbans/${kanban.kanban_id}`, { status_current: parseInt(kanban.status_current, 10) });
            if (dashboardType === "supplier") {
                setKanbans(prevKanbansByProduct => {
                    const updatedKanbansByProduct = { ...prevKanbansByProduct }; // Create a shallow copy
                    if (updatedKanbansByProduct[productID]) { // Check if productID exists
                        updatedKanbansByProduct[productID] = updatedKanbansByProduct[productID].map(k => {
                            if (k.kanban_id === kanban.kanban_id) {
                                return { ...k, status_name: response.data.status_name, status_color: response.data.status_color, customer_supplier: response.data.customer_supplier, status_current: response.data.status_current };
                            } else {
                                return k;
                            }
                        });
                    }
                    return updatedKanbansByProduct;
                });
            } else if (dashboardType === "customer") {
                setKanbans(prevKanbansByProduct => {
                    const updatedKanbansByProduct = { ...prevKanbansByProduct }; // Create a shallow copy
                    if (updatedKanbansByProduct[productID]) { // Check if productID exists
                        updatedKanbansByProduct[productID] = updatedKanbansByProduct[productID].map(k => {
                            if (k.kanban_id === kanban.kanban_id) {
                                return { ...k, status_name: response.data.status_name, status_color: response.data.status_color, customer_supplier: response.data.customer_supplier, status_current: response.data.status_current };
                            } else {
                                return k;
                            }
                        });
                    }
                    return updatedKanbansByProduct;
                });
            }
        } catch (error) {
            console.error("Error updating Kanban status", error);
        }
    };

    return (
        <div className="kanban-card" style={{ borderColor: kanban.status_color, borderWidth: '3px' }}>
            <h2>{kanban.status_name}</h2>  {/* Status name as main title, using <h2> */}
            <p><strong>Product:</strong> {kanban.product_name}</p>
            <p><strong>Container:</strong> {kanban.tipo_contenitore}</p>
            <p><strong>Quantity:</strong> {kanban.quantity}</p>
            <button
                style={{ backgroundColor: kanban.status_color, color: 'white', padding: '10px' }}
                onClick={handleStatusChange}
                disabled={dashboardType === 'supplier' ? kanban.customer_supplier !== 1 : kanban.customer_supplier !== 2}
            >
                Change Status
            </button>

        </div>
    );
};

export default KanbanCard;