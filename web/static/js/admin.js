// web/static/js/admin.js
document.addEventListener('DOMContentLoaded', function () {
    const uploadForm = document.getElementById('uploadForm');
    const uploadStatus = document.getElementById('uploadStatus');
    const dataPreview = document.getElementById('dataPreview');

    uploadForm.addEventListener('submit', async function (e) {
        e.preventDefault();

        const fileInput = document.getElementById('csvFile');
        const file = fileInput.files[0];
        if (!file) {
            showStatus('Selecteer eerst een CSV bestand', 'error');
            return;
        }

        const formData = new FormData();
        formData.append('file', file);

        try {
            const response = await fetch('/api/admin/upload', {
                method: 'POST',
                body: formData
            });

            const result = await response.json();

            if (response.ok) {
                showStatus(`Success! ${result.processed} records verwerkt.`, 'success');
                loadDataPreview();
            } else {
                showStatus(`Error: ${result.error}`, 'error');
            }
        } catch (error) {
            showStatus('Er is een fout opgetreden bij het uploaden.', 'error');
        }
    });

    function showStatus(message, type) {
        uploadStatus.textContent = message;
        uploadStatus.className = type;
    }

    async function loadDataPreview() {
        try {
            const response = await fetch('/api/admin/preview');
            const data = await response.json();

            if (data.length === 0) {
                dataPreview.innerHTML = '<p>Geen data beschikbaar</p>';
                return;
            }

            const table = document.createElement('table');
            table.innerHTML = `
                <tr>
                    <th>Studentnummer</th>
                    <th>Voornaam</th>
                    <th>Achternaam</th>
                    <th>Lokaal</th>
                    <th>Coach</th>
                </tr>
                ${data.map(student => `
                    <tr>
                        <td>${student.student_number}</td>
                        <td>${student.first_name}</td>
                        <td>${student.last_name}</td>
                        <td>${student.classroom}</td>
                        <td>${student.coach}</td>
                    </tr>
                `).join('')}
            `;

            dataPreview.innerHTML = '';
            dataPreview.appendChild(table);
        } catch (error) {
            dataPreview.innerHTML = '<p>Fout bij het laden van de preview</p>';
        }
    }

    // Load initial data preview
    loadDataPreview();
});