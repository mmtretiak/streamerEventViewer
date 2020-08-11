import {Component, OnInit, Pipe, PipeTransform} from '@angular/core';
import {StreamersService} from "../core/services/streamers.service";
import {DomSanitizer} from '@angular/platform-browser';
import {MatDialog} from "@angular/material/dialog";
import {AddStreamerDialog} from "./add-streamer/add-streamer.dialog";
import {ClipService} from "../core/services/clip.service";
import {NotificationService} from "../core/services/notification.service";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  activeStreamer: streamer;
  streamers: streamer[];

  totalViews: number;

  constructor(private streamersService: StreamersService, public dialog: MatDialog, private clipService: ClipService, private notificationService: NotificationService) { }
  ngOnInit(): void {
    this.getStreamers();
  }

  getStreamers() {
    this.streamersService.getStreamers().subscribe((res) => {
      if (res === null) {
        this.openDialog();
        return
      }

      this.clipService.getTotalViews().subscribe((res) => this.totalViews = res.count);
      this.clipService.getTotalViewsPerStreamer().subscribe((viewsRes) => {
        this.streamers = res as streamer[];
        this.streamers.map(x => {
          x.streamURL = 'https://player.twitch.tv/?channel=' + x.name + '&parent=stark-escarpment-52058.herokuapp.com&muted=true';
          x.chatURL = 'https://www.twitch.tv/embed/' + x.name + '/chat?parent=stark-escarpment-52058.herokuapp.com';
          if (viewsRes === null) {
            x.totalViews = 0;
          } else {
            const view = viewsRes.find(view => view.streamer_id == x.id);
            if (view !== undefined && view !== null) {
              x.totalViews = view.count;
            } else {
              x.totalViews = 0;
            }
          }
          return x;
        })
        this.activeStreamer = this.streamers[0];

        console.log(this.streamers);
      })
    });
  }

  changeStreamer(name: string) {
    this.activeStreamer = this.streamers.find(x => x.name == name);
  }

  openDialog() {
    const dialogRef = this.dialog.open(AddStreamerDialog, {
      width: '250px',
    });

    dialogRef.afterClosed().subscribe(result => {
      this.getStreamers();
    });
  }

  createClip() {
    this.clipService.addClip(this.activeStreamer.id).subscribe((res) => {
      this.notificationService.success("Clip added!");
      window.open(res.edit_url)
      console.log(res.edit_url)
    }, error => this.notificationService.warn(error));
  }
}

@Pipe({ name: 'safe' })
export class SafePipe implements PipeTransform {
  constructor(private sanitizer: DomSanitizer) {}
  transform(url) {
    return this.sanitizer.bypassSecurityTrustResourceUrl(url);
  }
}

interface streamer {
  id: string
  name: string;
  external_id: string;
  streamURL: string;
  chatURL: string;
  totalViews: number;
}
