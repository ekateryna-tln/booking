function CheckAvailability(csrf_token, room_id) {
    let html = `
        <form id="check-availability-form" action="" method="" novalidate class="needs-validation ps-3 pe-3">
            <div class="row" id="reservation-date-modal">
                <div class="col">
                    <div class="mb-3">
                        <label for="start_date" class="form-label">Starting date</label>
                        <input type="text" class="form-control" id="start_date" name="start_date"
                               aria-describedby="startDateHelp" required disabled autocomplete="off">
                    </div>
                </div>
                <div class="col">
                    <div class="mb-3">
                        <label for="end_date" class="form-label">Ending date</label>
                        <input type="text" class="form-control" id="end_date" name="end_date"
                               aria-describedby="endDateHelp" required disabled autocomplete="off">
                    </div>
                </div>
            </div>
        </form>
        `
    attention.custom({
        html: html,
        title: 'Choose your dates',
        willOpen: () => {
            const elem = document.getElementById('reservation-date-modal');
            const rangepicker = new DateRangePicker(elem, {
                format: "yyyy-mm-dd",
                showOnFocus: true,
                minDate: new Date()
            });
        },
        didOpen: () => {
            document.getElementById('start_date').removeAttribute('disabled');
            document.getElementById('end_date').removeAttribute('disabled');
        },
        callback: function(result) {

            let form = document.getElementById('check-availability-form')
            let formData = new FormData(form)
            formData.append('csrf_token', csrf_token)
            formData.append('room_id', room_id)

            fetch('/search-availability-json', {
                method: 'post',
                body: formData,
            })
                .then(response => response.json())
                .then(data => {
                    if (data.ok) {
                        attention.custom({
                            icon: 'success',
                            html: '<p>Room is available!</p>'
                                + '<p><a href="/book-room?uuid='
                                + data.room_id
                                + '&s='
                                + data.start_date
                                + '&e='
                                + data.end_date
                                + '" class="btn btn-primary">'
                                + 'Book now!</a></p>',
                            showConfirmButton: false
                        })
                    } else {
                        attention.error({
                            title: "No availability"
                        })
                    }
                })
        }
    })
}