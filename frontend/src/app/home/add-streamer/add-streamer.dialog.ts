import {Component, Inject} from "@angular/core";
import {StreamersService} from "../../core/services/streamers.service";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";

@Component({
  selector: 'add-streamer-dialog',
  templateUrl: 'add-streamer.dialog.html',
})
export class AddStreamerDialog {
  name: string;

  constructor(private streamersService: StreamersService,
              public dialogRef: MatDialogRef<AddStreamerDialog>) {
  }

  submit() {
    this.streamersService.addStreamer(this.name).subscribe(() => {
      this.dialogRef.close();
    })
  }
}
