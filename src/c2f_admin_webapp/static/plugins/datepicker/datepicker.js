'use strict';
$(document).ready(function() {
    $('.simple-date-range-picker').daterangepicker({
        // autoUpdateInput: true,
        startDate: moment().subtract(1, 'months'),
        ranges: {
            'Today': [moment(), moment()],
            'Yesterday' : [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
            'Last 7 Days' : [moment().subtract(6, 'days'), moment()],
            'Last 30 Days' : [moment().subtract(29, 'days'), moment()],
            'This Month': [moment().startOf('month'), moment().endOf('month')],
            'Last Month' : [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')]
        },
    });

    $('input[name="simple-date-range-picker-callback"]').daterangepicker({
        opens: 'left'
    }, function(start, end, label) {
        swal("A new date selection was made", start.format('DD-MM-YYYY') + ' to ' + end.format('DD-MM-YYYY'), "success")
    });
});