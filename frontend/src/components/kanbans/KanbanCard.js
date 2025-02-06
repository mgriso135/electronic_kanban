import React from 'react';
import api from '../../services/api';

const KanbanCard = ({ kanban, dashboardType, setKanbans, productID }) => {

    const handleStatusChange = async () => {
        console.log("KanbanCard - handleStatusChange CALLED - kanban_id:", kanban.kanban_id, "productID:", productID); // Log at start
        try {
            console.log("KanbanCard - handleStatusChange: Making API call..."); // Log before API call
            const response = await api.put(`/kanbans/${kanban.kanban_id}`, { status_current: parseInt(kanban.status_current, 10) });
            console.log("KanbanCard - handleStatusChange: API response received:", response.data); // Log API response

            if (dashboardType === "supplier" || dashboardType === "customer") { // Combined logic for both dashboards
                setKanbans(prevKanbansByProduct => {
                    console.log(`KanbanCard - setKanbans (${dashboardType} dashboard) CALLED with updatedKanban:`, response.data, `productID: (INSIDE CALLBACK - CHECK VALUE HERE):`, productID);

                    if (!prevKanbansByProduct[productID]) {
                        console.warn(`KanbanCard - setKanbans: productID "${productID}" not found in prevKanbansByProduct.`);
                        return prevKanbansByProduct || {}; // Return previous state to avoid errors
                    }


                    const updatedProductKanbans = prevKanbansByProduct[productID].map(k => {
                        if (k.kanban_id === kanban.kanban_id) {
                            return { ...k, status_name: response.data.status_name, status_color: response.data.status_color, customer_supplier: response.data.customer_supplier, status_current: response.data.status_current };
                        } else {
                            return k;
                        }
                    });

                    // Immutably update the state - create a new object and assign the updated product kanbans
                    const nextKanbansByProduct = {
                        ...prevKanbansByProduct,
                        [productID]: updatedProductKanbans,
                    };

                    return nextKanbansByProduct; // Return the updated state
                }, () => {
                    console.log(`KanbanCard - setKanbans (${dashboardType} dashboard) CALLBACK EXECUTED`);
                });
            }


        } catch (error) {
            console.error("KanbanCard - Error updating Kanban status", error);
        }
    };

    return (
        <div className="kanban-card" style={{ borderColor: kanban.status_color, borderWidth: '3px' }}>
            <h2>{kanban.status_name}</h2>
            <p><strong>Product:</strong> {kanban.product_name} (ID: {productID})</p>
            <p><strong>Container:</strong> {kanban.tipo_contenitore}</p>
            <p><strong>Quantity:</strong> {kanban.quantity}</p>
            {dashboardType === 'supplier' && kanban.customer_supplier === 1 || dashboardType === 'customer' && kanban.customer_supplier === 2 ? ( // Conditional rendering for button
                <button
                    style={{ backgroundColor: kanban.status_color, color: 'white', padding: '10px' }}
                    onClick={handleStatusChange}
                >
                    Change Status
                </button>
            ) : null} {/* Render nothing if not actionable */}
        </div>
    );
};

export default KanbanCard;