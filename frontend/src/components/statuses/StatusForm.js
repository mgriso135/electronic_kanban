import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const StatusForm = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [name, setName] = useState('');
    const [color, setColor] = useState('');
    const [isEdit, setIsEdit] = useState(false);

    useEffect(() => {
        if (id) {
            setIsEdit(true);
            const fetchStatus = async () => {
              try {
                  const response = await api.get(`/statuses/${id}`);
                  setName(response.data.name);
                  setColor(response.data.color);
              } catch(error) {
                console.error('Error fetching status:', error);
              }
            };
            fetchStatus();
        }
    }, [id]);

    const handleSubmit = async (e) => {
      e.preventDefault();
      const statusData = {
        name: name,
        color: color
      };
      try {
        if (isEdit) {
            await api.put(`/statuses/${id}`, statusData);
          } else {
              await api.post('/statuses', statusData)
          }
          navigate('/statuses');
      } catch (error) {
        console.error('Error saving status:', error);
      }
    };


    return (
      <div>
          <h2>{isEdit ? 'Edit Status' : 'Create New Status'}</h2>
            <form onSubmit={handleSubmit}>
                <div>
                  <label>Name:</label>
                  <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
                </div>
                <div>
                  <label>Color:</label>
                  <input type="text" value={color} onChange={(e) => setColor(e.target.value)} />
                </div>
                <button type="submit">{isEdit ? 'Update Status' : 'Create Status'}</button>
                <Link to="/statuses">Cancel</Link>
            </form>
        </div>
    );
};

export default StatusForm;